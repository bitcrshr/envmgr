
fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("compiling protos to rust...");

    let mut cfg = prost_build::Config::new();

    cfg.out_dir("./proto/compiled/rust")
        .compile_protos(&["./proto/src/environment_manager.proto"], &["./proto/src"])?;

    Ok(())
}
    
