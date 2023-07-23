use clap::{arg, command, Arg};

mod login;

fn main() {
    let matches = command!()
        .arg(Arg::new("login").exclusive(true))
        .get_matches();
        
    if matches.get_one::<String>("login").is_some() {
        println!("Logging in...");
        if let Err(e) = login::login() {
            println!("Error logging in: {}", e);
        }
    }
}