package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type submission struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Country      string    `json:"country"`
    SocialHandle string    `json:"socialHandle"`
    VideoURL     string    `json:"videoUrl"`
    Message      string    `json:"message"`
    CreatedAt    time.Time `json:"createdAt"`
}

type ballotNominee struct {
    ID           string    `json:"id"`
    SubmissionID string    `json:"submissionId"`
    Name         string    `json:"name"`
    Country      string    `json:"country"`
    SocialHandle string    `json:"socialHandle"`
    VideoURL     string    `json:"videoUrl"`
    Message      string    `json:"message"`
    Votes        int       `json:"votes"`
}

type ballotState struct {
    ID          string         `json:"id"`
    Title       string         `json:"title"`
    Description string         `json:"description"`
    Nominees    []ballotNominee `json:"nominees"`
    Active      bool           `json:"active"`
    CreatedAt   time.Time      `json:"createdAt"`
    ClosesAt    time.Time      `json:"closesAt"`
}

type fileStore struct {
    path        string
    mu          sync.Mutex
    submissions []submission
}

type contribution struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Amount    float64   `json:"amount"`
    Message   string    `json:"message"`
    CreatedAt time.Time `json:"createdAt"`
}

type contributionStore struct {
    path          string
    mu            sync.Mutex
    contributions []contribution
}

type ballotStore struct {
    path  string
    mu    sync.Mutex
    state ballotState
    votes map[string]string
}

func newFileStore(path string) (*fileStore, error) {
    store := &fileStore{path: path}
    if err := store.load(); err != nil {
        return nil, err
    }
    return store, nil
}

func (s *fileStore) load() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
        s.submissions = []submission{}
        return nil
    }

    data, err := os.ReadFile(s.path)
    if err != nil {
        return err
    }

    if len(data) == 0 {
        s.submissions = []submission{}
        return nil
    }

    if err := json.Unmarshal(data, &s.submissions); err != nil {
        return err
    }
    return nil
}

func (s *fileStore) save() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    data, err := json.MarshalIndent(s.submissions, "", "  ")
    if err != nil {
        return err
    }

    if err := os.WriteFile(s.path, data, 0o644); err != nil {
        return err
    }
    return nil
}

func (s *fileStore) add(sub submission) (submission, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    sub.ID = fmt.Sprintf("%d", time.Now().UnixNano())
    sub.CreatedAt = time.Now().UTC()
    s.submissions = append([]submission{sub}, s.submissions...)

    data, err := json.MarshalIndent(s.submissions, "", "  ")
    if err != nil {
        return submission{}, err
    }

    if err := os.WriteFile(s.path, data, 0o644); err != nil {
        return submission{}, err
    }

    return sub, nil
}

func (s *fileStore) list() []submission {
    s.mu.Lock()
    defer s.mu.Unlock()

    out := make([]submission, len(s.submissions))
    copy(out, s.submissions)
    return out
}

func newContributionStore(path string) (*contributionStore, error) {
    store := &contributionStore{path: path}
    if err := store.load(); err != nil {
        return nil, err
    }
    return store, nil
}

func (s *contributionStore) load() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
        s.contributions = []contribution{}
        return nil
    }

    data, err := os.ReadFile(s.path)
    if err != nil {
        return err
    }

    if len(data) == 0 {
        s.contributions = []contribution{}
        return nil
    }

    if err := json.Unmarshal(data, &s.contributions); err != nil {
        return err
    }
    return nil
}

func (s *contributionStore) add(entry contribution) (contribution, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    entry.ID = fmt.Sprintf("%d", time.Now().UnixNano())
    entry.CreatedAt = time.Now().UTC()
    s.contributions = append([]contribution{entry}, s.contributions...)

    data, err := json.MarshalIndent(s.contributions, "", "  ")
    if err != nil {
        return contribution{}, err
    }

    if err := os.WriteFile(s.path, data, 0o644); err != nil {
        return contribution{}, err
    }

    return entry, nil
}

func (s *contributionStore) list() []contribution {
    s.mu.Lock()
    defer s.mu.Unlock()

    out := make([]contribution, len(s.contributions))
    copy(out, s.contributions)
    return out
}

func (s *fileStore) getByID(id string) (submission, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()

    for _, sub := range s.submissions {
        if sub.ID == id {
            return sub, true
        }
    }
    return submission{}, false
}

func newBallotStore(path string) (*ballotStore, error) {
    bs := &ballotStore{
        path:  path,
        votes: make(map[string]string),
    }
    if err := bs.load(); err != nil {
        return nil, err
    }
    return bs, nil
}

func (b *ballotStore) load() error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if _, err := os.Stat(b.path); errors.Is(err, os.ErrNotExist) {
        b.state = ballotState{}
        b.votes = make(map[string]string)
        return nil
    }

    data, err := os.ReadFile(b.path)
    if err != nil {
        return err
    }

    if len(data) == 0 {
        b.state = ballotState{}
        b.votes = make(map[string]string)
        return nil
    }

    var wrapper struct {
        State ballotState        `json:"state"`
        Votes map[string]string `json:"votes"`
    }

    if err := json.Unmarshal(data, &wrapper); err != nil {
        return err
    }

    b.state = wrapper.State
    if wrapper.Votes != nil {
        b.votes = wrapper.Votes
    } else {
        b.votes = make(map[string]string)
    }
    return nil
}

func (b *ballotStore) save() error {
    b.mu.Lock()
    defer b.mu.Unlock()

    wrapper := struct {
        State ballotState        `json:"state"`
        Votes map[string]string `json:"votes"`
    }{
        State: b.state,
        Votes: b.votes,
    }

    data, err := json.MarshalIndent(wrapper, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(b.path, data, 0o644)
}

func (b *ballotStore) activeBallot() ballotState {
    b.mu.Lock()
    defer b.mu.Unlock()

    return b.state
}

func (b *ballotStore) setBallot(state ballotState) error {
    b.mu.Lock()
    b.state = state
    b.votes = make(map[string]string)
    b.mu.Unlock()
    return b.save()
}

func (b *ballotStore) addVote(email, nomineeID string) (ballotState, error) {
    email = strings.ToLower(strings.TrimSpace(email))
    if email == "" {
        return ballotState{}, errors.New("email required")
    }

    b.mu.Lock()
    defer b.mu.Unlock()

    if !b.state.Active {
        return ballotState{}, errors.New("no active ballot")
    }

    if _, exists := b.votes[email]; exists {
        return ballotState{}, errors.New("email already voted")
    }

    found := false
    for i := range b.state.Nominees {
        if b.state.Nominees[i].ID == nomineeID {
            b.state.Nominees[i].Votes++
            found = true
            break
        }
    }
    if !found {
        return ballotState{}, errors.New("nominee not found")
    }

    b.votes[email] = nomineeID

    if err := b.saveLocked(); err != nil {
        return ballotState{}, err
    }

    return b.state, nil
}

func (b *ballotStore) saveLocked() error {
    wrapper := struct {
        State ballotState        `json:"state"`
        Votes map[string]string `json:"votes"`
    }{
        State: b.state,
        Votes: b.votes,
    }

    data, err := json.MarshalIndent(wrapper, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(b.path, data, 0o644)
}

func main() {
    baseDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("failed to determine working directory: %v", err)
    }

    dataDir := filepath.Join(baseDir, "data")
    if err := os.MkdirAll(dataDir, 0o755); err != nil {
        log.Fatalf("failed to create data directory: %v", err)
    }

    store, err := newFileStore(filepath.Join(dataDir, "submissions.json"))
    if err != nil {
        log.Fatalf("failed to initialize store: %v", err)
    }

    ballotStore, err := newBallotStore(filepath.Join(dataDir, "ballot.json"))
    if err != nil {
        log.Fatalf("failed to initialize ballot store: %v", err)
    }

    bankStore, err := newContributionStore(filepath.Join(dataDir, "signal_bank.json"))
    if err != nil {
        log.Fatalf("failed to initialize contribution store: %v", err)
    }

    adminToken := strings.TrimSpace(os.Getenv("ORACLE_ADMIN_TOKEN"))

    mux := http.NewServeMux()

    mux.HandleFunc("/api/auditions", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            var payload struct {
                Name         string `json:"name"`
                Country      string `json:"country"`
                SocialHandle string `json:"socialHandle"`
                VideoURL     string `json:"videoUrl"`
                Message      string `json:"message"`
            }

            if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                http.Error(w, "invalid JSON", http.StatusBadRequest)
                return
            }

            payload.Name = strings.TrimSpace(payload.Name)
            payload.Country = strings.TrimSpace(payload.Country)
            payload.SocialHandle = strings.TrimSpace(payload.SocialHandle)
            payload.VideoURL = strings.TrimSpace(payload.VideoURL)
            payload.Message = strings.TrimSpace(payload.Message)

            if payload.Name == "" || payload.Country == "" || payload.VideoURL == "" {
                http.Error(w, "name, country, and videoUrl are required", http.StatusBadRequest)
                return
            }

            if !(strings.HasPrefix(payload.VideoURL, "http://") || strings.HasPrefix(payload.VideoURL, "https://")) {
                http.Error(w, "videoUrl must be a valid http(s) link", http.StatusBadRequest)
                return
            }

            created, err := store.add(submission{
                Name:         payload.Name,
                Country:      payload.Country,
                SocialHandle: payload.SocialHandle,
                VideoURL:     payload.VideoURL,
                Message:      payload.Message,
            })
            if err != nil {
                log.Printf("failed to store submission: %v", err)
                http.Error(w, "failed to record submission", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(created)

        case http.MethodGet:
            if adminToken != "" {
                token := r.URL.Query().Get("token")
                if token != adminToken {
                    http.Error(w, "unauthorized", http.StatusUnauthorized)
                    return
                }
            }

            submissions := store.list()

            countryFilter := strings.TrimSpace(r.URL.Query().Get("country"))
            if countryFilter != "" {
                filtered := make([]submission, 0, len(submissions))
                for _, sub := range submissions {
                    if strings.EqualFold(sub.Country, countryFilter) {
                        filtered = append(filtered, sub)
                    }
                }
                submissions = filtered
            }

            limitStr := strings.TrimSpace(r.URL.Query().Get("limit"))
            if limitStr != "" {
                if limit, err := strconv.Atoi(limitStr); err == nil && limit >= 0 && limit < len(submissions) {
                    submissions = submissions[:limit]
                }
            }

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(submissions)

        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/api/ballot", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            state := ballotStore.activeBallot()
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(state)
        case http.MethodPost:
            if adminToken != "" {
                token := r.URL.Query().Get("token")
                if token != adminToken {
                    http.Error(w, "unauthorized", http.StatusUnauthorized)
                    return
                }
            }

            var payload struct {
                Title       string   `json:"title"`
                Description string   `json:"description"`
                ClosesAt    string   `json:"closesAt"`
                NomineeIDs  []string `json:"nomineeIds"`
                Active      *bool    `json:"active"`
            }

            if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                http.Error(w, "invalid JSON", http.StatusBadRequest)
                return
            }

            if len(payload.NomineeIDs) == 0 {
                http.Error(w, "nomineeIds required", http.StatusBadRequest)
                return
            }

            nominees := make([]ballotNominee, 0, len(payload.NomineeIDs))
            for _, id := range payload.NomineeIDs {
                sub, ok := store.getByID(strings.TrimSpace(id))
                if !ok {
                    http.Error(w, fmt.Sprintf("submission %s not found", id), http.StatusBadRequest)
                    return
                }
                nominees = append(nominees, ballotNominee{
                    ID:           sub.ID,
                    SubmissionID: sub.ID,
                    Name:         sub.Name,
                    Country:      sub.Country,
                    SocialHandle: sub.SocialHandle,
                    VideoURL:     sub.VideoURL,
                    Message:      sub.Message,
                    Votes:        0,
                })
            }

            closesAt := time.Time{}
            if payload.ClosesAt != "" {
                ts, err := time.Parse(time.RFC3339, payload.ClosesAt)
                if err != nil {
                    http.Error(w, "closesAt must be RFC3339 timestamp", http.StatusBadRequest)
                    return
                }
                closesAt = ts
            }

            active := true
            if payload.Active != nil {
                active = *payload.Active
            }

            newState := ballotState{
                ID:          fmt.Sprintf("ballot-%d", time.Now().UnixNano()),
                Title:       strings.TrimSpace(payload.Title),
                Description: strings.TrimSpace(payload.Description),
                Nominees:    nominees,
                Active:      active,
                CreatedAt:   time.Now().UTC(),
                ClosesAt:    closesAt,
            }

            if err := ballotStore.setBallot(newState); err != nil {
                log.Printf("failed to set ballot: %v", err)
                http.Error(w, "failed to set ballot", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(newState)

        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/api/vote", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var payload struct {
            NomineeID string `json:"nomineeId"`
            Email     string `json:"email"`
            Name      string `json:"name"`
        }

        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "invalid JSON", http.StatusBadRequest)
            return
        }

        payload.NomineeID = strings.TrimSpace(payload.NomineeID)
        payload.Email = strings.TrimSpace(payload.Email)

        if payload.NomineeID == "" || payload.Email == "" {
            http.Error(w, "nomineeId and email are required", http.StatusBadRequest)
            return
        }

        if !strings.Contains(payload.Email, "@") {
            http.Error(w, "email must be valid", http.StatusBadRequest)
            return
        }

        updated, err := ballotStore.addVote(payload.Email, payload.NomineeID)
        if err != nil {
            switch err.Error() {
            case "email already voted":
                http.Error(w, err.Error(), http.StatusConflict)
                return
            case "no active ballot":
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            case "nominee not found":
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            default:
                log.Printf("failed to record vote: %v", err)
                http.Error(w, "failed to record vote", http.StatusInternalServerError)
                return
            }
        }

        votes := 0
        for _, nominee := range updated.Nominees {
            if nominee.ID == payload.NomineeID {
                votes = nominee.Votes
                break
            }
        }

        response := struct {
            BallotID  string `json:"ballotId"`
            NomineeID string `json:"nomineeId"`
            Votes     int    `json:"votes"`
            Message   string `json:"message"`
        }{
            BallotID:  updated.ID,
            NomineeID: payload.NomineeID,
            Votes:     votes,
            Message:   "Vote recorded. Thank you for supporting the contenders!",
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    mux.HandleFunc("/api/signal-bank/contributions", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            var payload struct {
                Name    string  `json:"name"`
                Amount  float64 `json:"amount"`
                Message string  `json:"message"`
            }

            if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                http.Error(w, "invalid JSON", http.StatusBadRequest)
                return
            }

            payload.Name = strings.TrimSpace(payload.Name)
            payload.Message = strings.TrimSpace(payload.Message)

            if payload.Amount <= 0 {
                http.Error(w, "amount must be greater than zero", http.StatusBadRequest)
                return
            }

            entry, err := bankStore.add(contribution{
                Name:    payload.Name,
                Amount:  payload.Amount,
                Message: payload.Message,
            })
            if err != nil {
                log.Printf("failed to store contribution: %v", err)
                http.Error(w, "failed to record contribution", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(entry)

        case http.MethodGet:
            contributions := bankStore.list()

            limitStr := strings.TrimSpace(r.URL.Query().Get("limit"))
            if limitStr != "" {
                if limit, err := strconv.Atoi(limitStr); err == nil && limit >= 0 && limit < len(contributions) {
                    contributions = contributions[:limit]
                }
            }

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(contributions)

        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    webDir := filepath.Join(baseDir, "web")
    mux.Handle("/", http.FileServer(http.Dir(webDir)))

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Digital Oracle server listening on :%s", port)
    if err := http.ListenAndServe(":"+port, loggingMiddleware(mux)); err != nil {
        log.Fatalf("server exited: %v", err)
    }
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
