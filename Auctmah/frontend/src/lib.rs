use wasm_bindgen::prelude::*;
use web_sys::{window, HtmlCanvasElement, CanvasRenderingContext2d, WebSocket};
use serde::{Deserialize, Serialize};
use std::cell::RefCell;
use std::rc::Rc;

#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(js_namespace = console)]
    pub fn log(s: &str);
}

// ============ DATA MODELS ============
#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Auction {
    pub id: String,
    pub title: String,
    pub description: String,
    pub start_price: f64,
    pub current_bid: f64,
    pub highest_bidder: String,
    pub bid_count: usize,
    pub status: String,
    pub start_time: String,
    pub end_time: String,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Bid {
    pub bidder_id: String,
    pub amount: f64,
    pub timestamp: String,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Message {
    #[serde(rename = "type")]
    pub msg_type: String,
    pub auction: Option<Auction>,
    pub bid: Option<Bid>,
}

// ============ AUCTION STATE ============
pub struct AuctionState {
    auctions: Vec<Auction>,
    selected_auction: Option<usize>,
    ws: Option<WebSocket>,
    canvas: Option<HtmlCanvasElement>,
    ctx: Option<CanvasRenderingContext2d>,
}

// ============ MAIN APP ============
#[wasm_bindgen(start)]
pub fn main() {
    log("🚀 Auctmah Rust+WASM Frontend Initializing...");
    
    let state = Rc::new(RefCell::new(AuctionState {
        auctions: vec![],
        selected_auction: None,
        ws: None,
        canvas: None,
        ctx: None,
    }));
    
    // Setup WebSocket connection
    connect_websocket(state.clone());
    
    // Setup canvas for rendering
    setup_canvas(state.clone());
    
    // Start animation loop
    render_loop(state.clone());
}

// ============ WEBSOCKET CONNECTION ============
fn connect_websocket(state: Rc<RefCell<AuctionState>>) {
    let window = window().expect("no global window");
    let location = window.location();
    let protocol = if location.protocol().unwrap_or_default() == "https:" {
        "wss"
    } else {
        "ws"
    };
    let host = location.host().unwrap_or_else(|_| "localhost:8080".to_string());
    let ws_url = format!("{}://{}/ws", protocol, host);
    
    match WebSocket::new(&ws_url) {
        Ok(ws) => {
            log(&format!("✅ Connected to {}", ws_url));
            
            let state_clone = state.clone();
            let onmessage = wasm_bindgen::closure::Closure::wrap(Box::new(move |event: web_sys::MessageEvent| {
                if let Some(msg_str) = event.data().as_string() {
                    if let Ok(msg) = serde_json::from_str::<Message>(&msg_str) {
                        handle_websocket_message(msg, state_clone.clone());
                    }
                }
            }) as Box<dyn FnMut(web_sys::MessageEvent)>);
            
            ws.set_onmessage(Some(onmessage.as_ref().unchecked_ref()));
            onmessage.forget();
            
            state.borrow_mut().ws = Some(ws);
        }
        Err(e) => {
            log(&format!("❌ WebSocket connection failed: {:?}", e));
        }
    }
}

// ============ MESSAGE HANDLER ============
fn handle_websocket_message(msg: Message, state: Rc<RefCell<AuctionState>>) {
    match msg.msg_type.as_str() {
        "auction_update" => {
            if let Some(auction) = msg.auction {
                let mut s = state.borrow_mut();
                if !s.auctions.iter().any(|a| a.id == auction.id) {
                    s.auctions.push(auction);
                    log(&format!("📢 New auction added (total: {})", s.auctions.len()));
                } else {
                    // Update existing auction
                    if let Some(existing) = s.auctions.iter_mut().find(|a| a.id == auction.id) {
                        *existing = auction;
                    }
                }
            }
        }
        "bid_accepted" => {
            if let Some(auction) = msg.auction {
                log(&format!("💰 Bid accepted! New highest: ${}", auction.current_bid));
                if let Some(existing) = state.borrow_mut().auctions.iter_mut().find(|a| a.id == auction.id) {
                    *existing = auction;
                }
            }
        }
        "auction_ended" => {
            if let Some(auction) = msg.auction {
                log(&format!("🏁 Auction ended: {}", auction.title));
                if let Some(existing) = state.borrow_mut().auctions.iter_mut().find(|a| a.id == auction.id) {
                    *existing = auction;
                }
            }
        }
        _ => {}
    }
}

// ============ CANVAS SETUP ============
fn setup_canvas(state: Rc<RefCell<AuctionState>>) {
    let window = window().expect("no global window");
    let document = window.document().expect("no document");
    
    match document.get_element_by_id("auction-canvas") {
        Some(canvas_element) => {
            if let Some(canvas) = canvas_element.dyn_ref::<HtmlCanvasElement>() {
                // Set canvas size
                canvas.set_width(1400);
                canvas.set_height(800);
                
                if let Ok(ctx) = canvas.get_context("2d") {
                    if let Some(ctx) = ctx.and_then(|ctx| ctx.dyn_into::<CanvasRenderingContext2d>().ok()) {
                        state.borrow_mut().canvas = Some(canvas.clone());
                        state.borrow_mut().ctx = Some(ctx);
                        log("✅ Canvas initialized (1400x800)");
                    }
                }
            }
        }
        None => {
            log("❌ Canvas element not found");
        }
    }
}

// ============ RENDER LOOP ============
fn render_loop(state: Rc<RefCell<AuctionState>>) {
    let state_clone = state.clone();
    
    let closure: Rc<RefCell<Option<wasm_bindgen::closure::Closure<dyn FnMut()>>>> = Rc::new(RefCell::new(None));
    let closure_clone = closure.clone();
    
    *closure.borrow_mut() = Some(wasm_bindgen::closure::Closure::wrap(Box::new(move || {
        render_frame(state_clone.clone());
        
        if let Some(window) = window() {
            window.request_animation_frame(
                closure_clone.borrow().as_ref().unwrap().as_ref().unchecked_ref()
            ).ok();
        }
    }) as Box<dyn FnMut()>));
    
    if let Some(window) = window() {
        window.request_animation_frame(
            closure.borrow().as_ref().unwrap().as_ref().unchecked_ref()
        ).ok();
    }
}

// ============ RENDER FRAME ============
fn render_frame(state: Rc<RefCell<AuctionState>>) {
    let state_ref = state.borrow();
    
    if let Some(ctx) = &state_ref.ctx {
        // Clear canvas
        ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#0a0e27"));
        ctx.fill_rect(0.0, 0.0, 1400.0, 800.0);
        
        // Draw grid
        draw_grid(ctx);
        
        // Draw auctions
        for (idx, auction) in state_ref.auctions.iter().enumerate() {
            let x = 50.0 + (idx % 3) as f64 * 450.0;
            let y = 50.0 + (idx / 3) as f64 * 350.0;
            
            draw_auction_card(ctx, auction, x, y);
        }
        
        // Draw header
        draw_header(ctx);
    }
}

// ============ DRAW HELPERS ============
fn draw_header(ctx: &CanvasRenderingContext2d) {
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#00d4ff"));
    ctx.set_font("32px Arial Bold");
    ctx.fill_text("🔨 AUCTMAH - Live Auction Board", 50.0, 40.0).ok();
    
    ctx.set_font("14px Arial");
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#a0aec0"));
    ctx.fill_text("Real-time bidding powered by Rust + WebAssembly", 50.0, 65.0).ok();
}

fn draw_grid(ctx: &CanvasRenderingContext2d) {
    ctx.set_stroke_style(&wasm_bindgen::JsValue::from_str("rgba(0, 212, 255, 0.1)"));
    ctx.set_line_width(1.0);
    
    // Vertical lines
    for i in 0..4 {
        let x = 50.0 + i as f64 * 450.0;
        ctx.begin_path();
        ctx.move_to(x, 100.0);
        ctx.line_to(x, 800.0);
        ctx.stroke();
    }
    
    // Horizontal lines
    for i in 0..3 {
        let y = 50.0 + i as f64 * 350.0;
        ctx.begin_path();
        ctx.move_to(0.0, y);
        ctx.line_to(1400.0, y);
        ctx.stroke();
    }
}

fn draw_auction_card(ctx: &CanvasRenderingContext2d, auction: &Auction, x: f64, y: f64) {
    let w = 400.0;
    let h = 280.0;
    
    // Background
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("rgba(17, 22, 51, 0.8)"));
    ctx.fill_rect(x, y, w, h);
    
    // Border
    let border_color = match auction.status.as_str() {
        "active" => "#00d4ff",
        "ended" => "#ff6b6b",
        _ => "#8b5cf6",
    };
    ctx.set_stroke_style(&wasm_bindgen::JsValue::from_str(border_color));
    ctx.set_line_width(2.0);
    ctx.stroke_rect(x, y, w, h);
    
    // Title
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#00d4ff"));
    ctx.set_font("16px Arial Bold");
    ctx.fill_text(&auction.title, x + 15.0, y + 30.0).ok();
    
    // Status badge
    let status_bg = match auction.status.as_str() {
        "active" => "#00d4ff",
        "ended" => "#ff6b6b",
        _ => "#8b5cf6",
    };
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str(status_bg));
    ctx.fill_rect(x + w - 110.0, y + 10.0, 100.0, 25.0);
    
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#0a0e27"));
    ctx.set_font("12px Arial Bold");
    ctx.fill_text(&auction.status.to_uppercase(), x + w - 95.0, y + 27.0).ok();
    
    // Current bid (LARGE)
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#ff6b6b"));
    ctx.set_font("28px Arial Bold");
    ctx.fill_text(&format!("${:.0}", auction.current_bid), x + 15.0, y + 80.0).ok();
    
    // Bid count
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#a0aec0"));
    ctx.set_font("12px Arial");
    ctx.fill_text(&format!("Bids: {}", auction.bid_count), x + 15.0, y + 110.0).ok();
    
    // Highest bidder
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#a0aec0"));
    ctx.set_font("12px Arial");
    ctx.fill_text(&format!("Leader: {}", &auction.highest_bidder[..std::cmp::min(15, auction.highest_bidder.len())]), x + 15.0, y + 130.0).ok();
    
    // Progress bar (time remaining)
    draw_progress_bar(ctx, x + 15.0, y + 170.0, 370.0, 15.0);
    
    // Description
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#a0aec0"));
    ctx.set_font("11px Arial");
    let desc = if auction.description.len() > 40 {
        &auction.description[..40]
    } else {
        &auction.description
    };
    ctx.fill_text(desc, x + 15.0, y + 220.0).ok();
    
    // Bottom action
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("rgba(0, 212, 255, 0.2)"));
    ctx.fill_rect(x + 15.0, y + 240.0, 370.0, 30.0);
    
    ctx.set_fill_style(&wasm_bindgen::JsValue::from_str("#00d4ff"));
    ctx.set_font("14px Arial Bold");
    ctx.fill_text("Click to bid →", x + 160.0, y + 260.0).ok();
}fn draw_progress_bar(ctx: &CanvasRenderingContext2d, x: f64, y: f64, w: f64, h: f64) {
    // Background
    ctx.set_fill_style(&JsValue::from_str("rgba(255, 107, 107, 0.2)"));
    ctx.fill_rect(x, y, w, h);
    
    // Progress (animate)
    let progress = ((js_sys::Date::now() as u32 % 10000) as f64 / 10000.0) * 100.0;
    ctx.set_fill_style(&JsValue::from_str("#ff6b6b"));
    ctx.fill_rect(x, y, w * (progress / 100.0), h);
}

// ============ PLACE BID FUNCTION ============
#[wasm_bindgen]
pub fn place_bid(auction_id: String, amount: f64) {
    log(&format!("💰 Placing bid on {} for ${}", auction_id, amount));
    
    if let Some(_window) = window() {
        if let Some(ws) = get_websocket() {
            let bid = Bid {
                bidder_id: format!("bidder_{}", js_sys::Date::now() as u32),
                amount,
                timestamp: js_sys::Date::now().to_string(),
            };
            
            let msg = Message {
                msg_type: "place_bid".to_string(),
                auction: None,
                bid: Some(bid),
            };
            
            if let Ok(msg_str) = serde_json::to_string(&msg) {
                ws.send_with_str(&msg_str).ok();
            }
        }
    }
}

fn get_websocket() -> Option<WebSocket> {
    // This is a placeholder - in real implementation, store WS in global
    None
}
