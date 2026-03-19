package deployer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"stackd/internal/cli"
	"strings"
)

func HandlePythonSetup(absPath string, reader *bufio.Reader) string {
	venvPaths := []string{".venv", "venv", "env"}
	venvFound := ""
	for _, v := range venvPaths {
		if fileExists(filepath.Join(absPath, v)) {
			venvFound = v
			break
		}
	}

	if venvFound == "" {
		fmt.Println("\n" + cli.Yellow + "⚠️  Perhatian: Folder Virtual Environment (venv) tidak terdeteksi." + cli.Reset)
		fmt.Println("Aplikasi Python disarankan berjalan terisolasi menggunakan venv.")
		fmt.Print(cli.Bold + "➤ Automasi membuat venv & install requirements.txt sekarang? (y/n) [y]: " + cli.Reset)

		createVenv, _ := reader.ReadString('\n')
		createVenv = strings.TrimSpace(strings.ToLower(createVenv))

		if createVenv == "" || createVenv == "y" {
			fmt.Println(cli.Blue + "⏳ [1/2] Membuat virtual environment (python3 -m venv venv)..." + cli.Reset)

			err := runSystemCommandInDir(absPath, "python3", "-m", "venv", "venv")
			if err != nil {
				fmt.Printf(cli.Red+"❌ Gagal membuat venv: %v\n"+cli.Reset, err)
				printPythonManualGuide(absPath)
				os.Exit(1)
			}
			venvFound = "venv"

			if fileExists(filepath.Join(absPath, "requirements.txt")) {
				fmt.Println(cli.Blue + "⏳ [2/2] Menginstall dari requirements.txt..." + cli.Reset)
				pipCmd := filepath.Join(absPath, "venv", "bin", "pip")
				err = runSystemCommandInDir(absPath, pipCmd, "install", "-r", "requirements.txt")
				if err != nil {
					fmt.Printf(cli.Red+"❌ Peringatan: Gagal menginstall requirements secara otomatis: %v\n"+cli.Reset, err)
				} else {
					fmt.Println(cli.Green + "✅ Venv & Dependensi berhasil diinstall." + cli.Reset)
				}
			} else {
				fmt.Println(cli.Green + "✅ Venv berhasil dibuat! (file requirements.txt tidak ditemukan)." + cli.Reset)
			}
		} else {
			printPythonManualGuide(absPath)
			fmt.Print("\n" + cli.Bold + "Lanjutkan ke konfigurasi systemd tanpa venv? (y/n) [n]: " + cli.Reset)
			cont, _ := reader.ReadString('\n')
			cont = strings.TrimSpace(strings.ToLower(cont))
			if cont != "y" {
				fmt.Println(cli.Red + "Dibatalkan." + cli.Reset)
				os.Exit(0)
			}
		}
	} else {
		fmt.Printf(cli.Green+"✅ Virtual Environment Python terdeteksi di: %s\n"+cli.Reset, venvFound)
	}

	if venvFound != "" {
		gunicornPath := filepath.Join(absPath, venvFound, "bin", "gunicorn")
		pythonPath := filepath.Join(absPath, venvFound, "bin", "python")

		fmt.Print("\n" + cli.Bold + "➤ Apakah Anda menggunakan Gunicorn/Uvicorn untuk production? (y/t) [y]: " + cli.Reset)
		useWgsi, _ := reader.ReadString('\n')
		useWgsi = strings.TrimSpace(strings.ToLower(useWgsi))

		if useWgsi == "" || useWgsi == "y" || useWgsi == "ya" {
			return fmt.Sprintf("%s -w 4 -b 127.0.0.1:5000 app:app", gunicornPath)
		}
		return fmt.Sprintf("%s app.py", pythonPath)
	}

	return "/usr/bin/python3 app.py"
}

func printPythonManualGuide(absPath string) {
	fmt.Println("\n" + cli.Yellow + "💡 Panduan manual setup Python Production:" + cli.Reset)
	fmt.Println("  1. Masuk ke direktori: cd", absPath)
	fmt.Println("  2. Buat venv         : python3 -m venv venv")
	fmt.Println("  3. Aktifkan venv     : source venv/bin/activate")
	fmt.Println("  4. Install library   : pip install -r requirements.txt")
	fmt.Println("  5. Install gunicorn  : pip install gunicorn")
}

func HandleNodeSetup(absPath string, reader *bufio.Reader) string {
	if !fileExists(filepath.Join(absPath, "node_modules")) && fileExists(filepath.Join(absPath, "package.json")) {
		fmt.Println("\n" + cli.Yellow + "⚠️  Folder node_modules tidak ditemukan." + cli.Reset)
		fmt.Print(cli.Bold + "➤ Apakah Anda ingin menjalankan 'npm install' secara otomatis? (y/n) [y]: " + cli.Reset)
		npmInst, _ := reader.ReadString('\n')
		npmInst = strings.TrimSpace(strings.ToLower(npmInst))
		if npmInst == "" || npmInst == "y" {
			fmt.Println(cli.Blue + "⏳ Menjalankan npm install..." + cli.Reset)
			err := runSystemCommandInDir(absPath, "npm", "install")
			if err != nil {
				fmt.Printf(cli.Red+"❌ Peringatan: npm install gagal: %v\n"+cli.Reset, err)
			} else {
				fmt.Println(cli.Green + "✅ npm install berhasil." + cli.Reset)
			}
		}
	}

	return "/usr/bin/npm run start"
}

func HandleLaravelSetup(absPath string, reader *bufio.Reader) string {
	if !fileExists(filepath.Join(absPath, "vendor")) && fileExists(filepath.Join(absPath, "composer.json")) {
		fmt.Println("\n" + cli.Yellow + "⚠️  Folder vendor tidak ditemukan." + cli.Reset)
		fmt.Print(cli.Bold + "➤ Apakah Anda ingin stackd menjalankan 'composer install'? (y/n) [y]: " + cli.Reset)
		cInst, _ := reader.ReadString('\n')
		cInst = strings.TrimSpace(strings.ToLower(cInst))
		if cInst == "" || cInst == "y" {
			fmt.Println(cli.Blue + "⏳ Menjalankan composer install --no-dev..." + cli.Reset)
			err := runSystemCommandInDir(absPath, "composer", "install", "--no-dev", "--optimize-autoloader")
			if err != nil {
				fmt.Printf(cli.Red+"❌ Peringatan: composer install gagal (pastikan composer terinstall di sistem): %v\n"+cli.Reset, err)
			} else {
				fmt.Println(cli.Green + "✅ composer install berhasil." + cli.Reset)
			}
		}
	}

	return "/usr/bin/php artisan serve --host=0.0.0.0 --port=8000"
}

func runSystemCommandInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
