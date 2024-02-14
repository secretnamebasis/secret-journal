
# secret-journal

`secret-journal` is a GUI-based encrypted journal powered by [DERO](https://github.com/deroproject/derohe) homomorphically-encrypted, blockchain [technology](https://github.com/deroproject/graviton) and provides a secure and private space for users to log their thoughts, ideas, and experiences.


## Configurations

**Login/Logout**: In order to login to secret-journal, you will need to have a DERO wallet running in RPC Server Mode: 
- [CLI](https://github.com/deroproject/derohe/releases/latest)
- [GUI](https://github.com/DEROFDN/Engram/releases/latest)

**Balance**: In order to make entries into you DERO wallet, you will need to have a balance of DERO for new entry fees.  

## Features

- **Latest Entries**: Quickly access the most recent journal entries for a seamless experience.
  
- **Privacy Button**: Toggle the visibility of journal contents with the privacy button, ensuring enhanced confidentiality.

- **Search Field**: Easily find specific entries using the search functionality.

- **Chat-style**: Create new journal entries effortlessly, capturing your thoughts in real-time.

## Getting Started

### Latest Release

Grab the latest [release](https://github.com/secretnamebasis/secret-journal/releases/latest/).

> [android](https://github.com/secretnamebasis/secret-journal/releases/latest/secret-journal.apk) & [linux](https://github.com/secretnamebasis/secret-journal/releases/latest/secret-journal.tar.xz) 


### From Source

1. Clone the repository:

   ```bash
   git clone https://github.com/secretnamebasis/secret-journal.git
   cd secret-journal
   ```

2. Build the application (make sure you have Go installed):

   ```bash
   go build .
   ```

3. Run the compiled binary:

   ```bash
   ./secret-journal
   ```


## Usage

1. Launch the Secret Journal application.
2. Log in to your DERO wallet RPC to access your private journal.
3. Explore the intuitive GUI with features like the latest entries, privacy button, search field, and new entry creation.

## REVENUE NOTICE: 
 `secret-journal` has a small tip of `0.002 DERO` included with new entries by default to support [secretsystems](https://github.com/secretsystems) development. 

## Contributing

If you would like to contribute to Secret Journal, please follow your heart and make a PR.

## License

This project is licensed under the [GNU PLv3](LICENSE).
