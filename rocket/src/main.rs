#![feature(proc_macro_hygiene)]
#[macro_use] extern crate rocket;


#[get("/")]
fn hello() -> &'static str {
    "Hello, world2!"
}

fn main() {
    let _ = rocket::ignite().mount("/", routes![hello]).launch();
}
