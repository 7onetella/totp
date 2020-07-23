package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	"os"

	imgcat "github.com/martinlindhe/imgcat/lib"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/xlzd/gotp"
	"rsc.io/qr"
)

var showQRCode bool
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "totp <issuer>",
	Short: "calculate TOTP for given issuer",
	Long:  `calculate TOTP for given issuer`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("missing issuer argument")
			return
		}

		issuerName := args[0]

		home := getTOTPHome()

		fileName := fmt.Sprintf("%s/%s.txt", home, issuerName)

		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			terms := strings.Split(line, " ")
			account := terms[0]
			secret := strings.TrimSpace(terms[1])
			totp := gotp.NewDefaultTOTP(secret)

			fmt.Println()
			fmt.Println("account    :", account)
			if showQRCode {
				renderQRCode(totp, account, issuerName)
				continue
			}

			fmt.Println(totp.Now())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flags := rootCmd.Flags()
	flags.BoolVar(&showQRCode, "qr", false, "show qr code")
	flags.BoolVar(&debug, "debug", false, "debug")
}

func getTOTPHome() string {
	homedir, _ := homedir.Dir() // return path with slash at the end
	totpHome := homedir + "/.totp"

	return totpHome
}

func generateQRCodeImage(url string) ([]byte, error) {
	code, err := qr.Encode(url, qr.Q)
	if err != nil {
		return nil, err
	}
	return code.PNG(), nil
}

func renderQRCode(totp *gotp.TOTP, accountName, issuerName string) {
	if len(accountName) == 0 {
		fmt.Println("you must specify --account flag value")
		return
	}

	url := totp.ProvisioningUri(accountName, issuerName)

	if debug {
		fmt.Println("QR encoded : " + url)
	}

	qrBytes, _ := generateQRCodeImage(url)

	r := bytes.NewReader(qrBytes)
	imgcat.Cat(r, os.Stdout)
	return
}

func main() {
	Execute()
}
