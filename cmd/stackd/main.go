package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"stackd/internal/cli"
	"stackd/internal/deployer"
	"stackd/internal/detector"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "help", "-h", "--help":
		printHelp()
		os.Exit(0)
	case "deploy":
		if len(os.Args) < 3 {
			fmt.Println(cli.Red + "❌ Error: Path project tidak disertakan." + cli.Reset)
			fmt.Println(cli.Yellow + "💡 Penggunaan: stackd deploy <path>" + cli.Reset)
			os.Exit(1)
		}
		path := os.Args[2]
		handleDeploy(path)
	case "remove", "undeploy":
		if len(os.Args) < 3 {
			fmt.Println(cli.Red + "❌ Error: Nama service tidak disertakan." + cli.Reset)
			fmt.Println(cli.Yellow + "💡 Penggunaan: stackd remove <nama_service>" + cli.Reset)
			os.Exit(1)
		}
		serviceName := os.Args[2]
		handleRemove(serviceName)
	case "list":
		deployer.ListManagedServices()
	default:
		fmt.Printf(cli.Red+"❌ Perintah tidak dikenal: %s\n"+cli.Reset, command)
		os.Exit(1)
	}
}

func handleDeploy(relPath string) {
	fmt.Println(cli.Cyan + cli.Bold + "=====================================" + cli.Reset)
	fmt.Println(cli.Cyan + cli.Bold + "🚀 stackd CLI - Deployment Tools" + cli.Reset)
	fmt.Println(cli.Cyan + cli.Bold + "=====================================" + cli.Reset)

	absPath, err := filepath.Abs(relPath)
	if err != nil {
		fmt.Printf(cli.Red+"❌ Error membaca path: %v\n"+cli.Reset, err)
		os.Exit(1)
	}

	if !fileExists(absPath) {
		fmt.Printf(cli.Red+"❌ Error: Path tidak ditemukan (%s)\n"+cli.Reset, absPath)
		os.Exit(1)
	}

	fmt.Println(cli.Blue + "🔍 Menganalisa project..." + cli.Reset)
	detectedType := detector.DetectProject(absPath)
	fmt.Printf(cli.Green+"✅ Terdeteksi sebagai  : %s\n"+cli.Reset, strings.ToUpper(detectedType))

	reader := bufio.NewReader(os.Stdin)

	// Pilihan Tipe Project
	fmt.Printf("\n"+cli.Bold+"➤ Pilih tipe project (misal: python, node, laravel, go) [%s]: "+cli.Reset, detectedType)
	userType, _ := reader.ReadString('\n')
	userType = strings.TrimSpace(userType)
	if userType == "" {
		userType = detectedType
	}

	// Nama Service
	defaultServiceName := filepath.Base(absPath)
	fmt.Printf(cli.Bold+"➤ Nama service systemd (misal: %s_app) [%s]: "+cli.Reset, defaultServiceName, defaultServiceName)
	serviceName, _ := reader.ReadString('\n')
	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		serviceName = defaultServiceName
	}

	// Command Eksekusi
	defaultCmd := ""
	switch userType {
	case "python":
		defaultCmd = deployer.HandlePythonSetup(absPath, reader)
	case "node":
		defaultCmd = deployer.HandleNodeSetup(absPath, reader)
	case "laravel":
		defaultCmd = deployer.HandleLaravelSetup(absPath, reader)
	case "go":
		defaultCmd = filepath.Join(absPath, serviceName)
	default:
		defaultCmd = "./run.sh"
	}

	fmt.Printf("\n" + cli.Yellow + "⚠️  Command untuk run program sangat penting (disarankan gunakan absolute path)." + cli.Reset + "\n")
	fmt.Printf(cli.Yellow + "Contoh Node  : /usr/bin/npm run start" + cli.Reset + "\n")
	fmt.Printf(cli.Bold+"➤ Command run program [%s]: "+cli.Reset, defaultCmd)
	runCmd, _ := reader.ReadString('\n')
	runCmd = strings.TrimSpace(runCmd)
	if runCmd == "" {
		runCmd = defaultCmd
	}

	// User Eksekusi
	defaultUser := "root"
	fmt.Printf(cli.Bold+"➤ Jalankan service sebagai user linux apa? (misal: www-data, ubuntu) [%s]: "+cli.Reset, defaultUser)
	execUser, _ := reader.ReadString('\n')
	execUser = strings.TrimSpace(execUser)
	if execUser == "" {
		execUser = defaultUser
	}

	fmt.Println("\n" + cli.Cyan + "=====================================" + cli.Reset)
	fmt.Println(cli.Bold + "📋 Ringkasan Deployment" + cli.Reset)
	fmt.Println(cli.Cyan + "=====================================" + cli.Reset)
	fmt.Printf("Path         : %s\n", absPath)
	fmt.Printf("Tipe         : %s\n", userType)
	fmt.Printf("Service Name : %s\n", serviceName)
	fmt.Printf("Command      : %s\n", runCmd)
	fmt.Printf("Run As User  : %s\n", execUser)
	fmt.Print("\n" + cli.Bold + "Lanjutkan membuat systemd service? (y/n) [y]: " + cli.Reset)

	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" && confirm != "" {
		fmt.Println(cli.Red + "🛑 Dibatalkan oleh user." + cli.Reset)
		os.Exit(0)
	}

	fmt.Println("\n" + cli.Blue + "⚙️  Menyiapkan systemd..." + cli.Reset)
	err = deployer.CreateSystemdService(serviceName, absPath, runCmd, execUser)
	if err != nil {
		fmt.Printf("\n"+cli.Red+"❌ Gagal membuat service:\n%v\n\n"+cli.Yellow+"💡 Petunjuk: Apakah kamu sudah menjalankannya dengan 'sudo'? File systemd membutuhkan akses root.\n"+cli.Reset, err)
		os.Exit(1)
	}

	fmt.Println("\n" + cli.Green + cli.Bold + "🎉 Deployment systemd service berhasil!" + cli.Reset)
	fmt.Println(cli.Cyan + "-------------------------------------" + cli.Reset)
	fmt.Printf("Cek status  : sudo systemctl status %s\n", serviceName)
	fmt.Printf("Lihat log   : sudo journalctl -u %s -f\n", serviceName)
	fmt.Printf("Restart     : sudo systemctl restart %s\n", serviceName)
}

func handleRemove(serviceName string) {
	fmt.Printf(cli.Yellow+"⚠️  Anda yakin ingin menghapus service '%s' secara permanen dari server? (y/n) [n]: "+cli.Reset, serviceName)
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" {
		fmt.Println(cli.Red + "🛑 Penghapusan dibatalkan." + cli.Reset)
		os.Exit(0)
	}

	err := deployer.RemoveSystemdService(serviceName)
	if err != nil {
		fmt.Printf("\n"+cli.Red+"❌ Gagal menghapus service:\n%v\n\n"+cli.Yellow+"💡 Petunjuk: Apakah Anda lupa menggunakan hak 'sudo'? Menghapus service memerlukan akses root.\n"+cli.Reset, err)
		os.Exit(1)
	}
	fmt.Printf("\n"+cli.Green+"✅ Seluruh konfigurasi service '%s' telah bersih dihapus dari sistem!\n"+cli.Reset, serviceName)
}

func printHelp() {
	fmt.Println(cli.Cyan + cli.Bold + "stackd" + cli.Reset + " - Ultimate CLI Deployment Automator")
	fmt.Println("\n" + cli.Bold + "PENGGUNAAN:" + cli.Reset)
	fmt.Println("  stackd <perintah> [argumen]")
	fmt.Println("\n" + cli.Bold + "PERINTAH TERSEDIA:" + cli.Reset)
	fmt.Println("  " + cli.Green + "deploy" + cli.Reset + " <path>   Mendeploy project ke systemd (Gunakan absolute/relative path)")
	fmt.Println("  " + cli.Red + "remove" + cli.Reset + " <nama>   Menghapus/unregister sebuah service dari systemd")
	fmt.Println("  " + cli.Blue + "list" + cli.Reset + "            List service/aplikasi user deploy yang aktif di systemd")
	fmt.Println("  " + cli.Yellow + "help, -h" + cli.Reset + "        Menampilkan layar bantuan ini")
	fmt.Println("\n" + cli.Bold + "CONTOH:" + cli.Reset)
	fmt.Println("  sudo stackd deploy /var/www/ecommerce-api")
	fmt.Println("  sudo stackd remove ecommerce-api")
	fmt.Println("  stackd list")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
