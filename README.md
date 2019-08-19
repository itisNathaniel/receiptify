# receiptify
#### Import email receipts to Monzo 

---

#### How to get started

1. Install Go! [Official Docs](https://golang.org/doc/install) and [DigitalOcean Guide](https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-macos) might help.
1. Clone the repo into your Go `src` directory
3. Rename `.example.env` to `.env`, fill in the environment variables. `accesstoken` is your access token from [developer.monzo.com](//developer.monzo.com)
4. `cd` into your the receiptify directory
5. Import your environment variables using `set -a`, `source .env` and `set +a`
6. Run it! Use `go run *.go` 

---

#### 🚧 Your milage may vary

- This works with reciepts from Trainline I've tested from 2017-present and with Wetherspoons reciepts since 2017
- If you fix any bugs or add new merchants, please feel free to PR them in!
- Monzo's API returns 403 for prepaid account matches, so we can't match those currently. See more [in the post](//nathaniel.work/receiptify-post)
