package main

import (
	"bytes"
	"fmt"
	"strings"

	"io/ioutil"
	"os"

	imgcat "github.com/martinlindhe/imgcat/lib"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/xlzd/gotp"
	"rsc.io/qr"
)

var showQRCode bool
var accountName string
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

		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.txt", home, issuerName))
		if err != nil {
			fmt.Println(err)
			return
		}
		secret := string(data)
		secret = strings.TrimSpace(secret)

		totp := gotp.NewDefaultTOTP(secret)

		if showQRCode {
			if len(accountName) == 0 {
				fmt.Println("you must specify --account flag value")
				return
			}

			url := totp.ProvisioningUri(accountName, issuerName)
			if debug {
				fmt.Println()
				fmt.Println("encoded url = " + url)
			}
			qrBytes, _ := generateQRCodeImage(url)

			fmt.Println()
			r := bytes.NewReader(qrBytes)
			imgcat.Cat(r, os.Stdout)
			return
		}

		fmt.Println(totp.Now())
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
	flags.StringVar(&accountName, "account", "", "account name")
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

func main() {
	Execute()
}
