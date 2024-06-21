use std::{io::{Read, Write}, net::TcpStream};

use x25519_dalek::{EphemeralSecret, PublicKey};
use rand_chacha::rand_core::SeedableRng;

const SERVER_PUB_KEY: [u8; 32] = [4, 113, 76, 112, 152, 193, 181, 197, 53, 136, 194, 15, 205, 130, 160, 81, 65, 93, 34, 247, 245, 169, 24, 115, 25, 176, 214, 72, 142, 1, 224, 47];

// fn gen_key() -> ([u8], [u8]) {

//     let rng = rand_chacha::ChaCha20Rng::from_entropy();
//     let bob_secret = EphemeralSecret::random_from_rng(rng);
//     let bob_public = PublicKey::from(&bob_secret);

//     let alice_shared_secret = alice_secret.diffie_hellman(&bob_public);
//     let bob_shared_secret = bob_secret.diffie_hellman(&alice_public);
//     assert_eq!(alice_shared_secret.as_bytes(), bob_shared_secret.as_bytes());
//     println!("{:?}", alice_shared_secret.as_bytes());
//     println!("Hello, world!");
// }

fn main() -> std::io::Result<()>{
    // Generate keys
    let rng = rand_chacha::ChaCha20Rng::from_entropy();
    let client_secret = EphemeralSecret::random_from_rng(rng);
    let client_public = PublicKey::from(&client_secret);
    let server_public = PublicKey::from(SERVER_PUB_KEY);

    // Generate share secret
    let shared_secret = client_secret.diffie_hellman(&server_public);

    // Encrypt message with shared secret
    println!("shared secret: {:?}", shared_secret.to_bytes());


    let mut stream = TcpStream::connect("127.0.0.1:9090")?;
    let msg = b"hello world\n".as_slice();

    // Write public key
    let _ = stream.write(client_public.as_bytes());
    println!("request: {:?}", client_public.as_bytes());

    // Write message
    stream.write(msg)?;
    println!("request: {:?}", msg);

    // Decrypt response with server public key
    let mut buf = [0; 128];
    let size = stream.read(&mut buf)?;
    println!("response: {:?}", &buf[0..size]);

    Ok(())
}
