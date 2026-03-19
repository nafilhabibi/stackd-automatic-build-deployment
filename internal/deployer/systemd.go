package deployer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"stackd/internal/cli"
	"strings"
)

func CreateSystemdService(serviceName, workDir, execStart, user string) error {
	serviceContent := fmt.Sprintf(`[Unit]
Description=%s Service
After=network.target

[Service]
Type=simple
User=%s
WorkingDirectory=%s
ExecStart=%s
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
# Managed by stackd
`, serviceName, user, workDir, execStart)

	serviceFilePath := filepath.Join("/etc/systemd/system", fmt.Sprintf("%s.service", serviceName))

	// Write the service file
	err := os.WriteFile(serviceFilePath, []byte(serviceContent), 0644)
	if err != nil {
		return fmt.Errorf("Gagal menulis file unit systemd ke %s: %v", serviceFilePath, err)
	}

	fmt.Printf(cli.Green+"✅ File service telah dibuat di: %s\n"+cli.Reset, serviceFilePath)

	if err := runSystemCommand("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("Gagal reload daemon: %v", err)
	}

	if err := runSystemCommand("systemctl", "enable", fmt.Sprintf("%s.service", serviceName)); err != nil {
		return fmt.Errorf("Gagal enable service: %v", err)
	}

	if err := runSystemCommand("systemctl", "start", fmt.Sprintf("%s.service", serviceName)); err != nil {
		return fmt.Errorf("Gagal start service: %v", err)
	}

	return nil
}

func RemoveSystemdService(serviceName string) error {
	serviceFileName := fmt.Sprintf("%s.service", serviceName)
	serviceFilePath := filepath.Join("/etc/systemd/system", serviceFileName)

	if !fileExists(serviceFilePath) {
		return fmt.Errorf("File konfigurasi untuk service '%s' tidak ditemukan di %s", serviceName, serviceFilePath)
	}

	fmt.Printf(cli.Yellow+"🛑 Menghentikan service %s...\n"+cli.Reset, serviceName)
	_ = runSystemCommand("systemctl", "stop", serviceFileName)

	fmt.Printf(cli.Yellow+"🛑 Menonaktifkan service %s dari autostart...\n"+cli.Reset, serviceName)
	_ = runSystemCommand("systemctl", "disable", serviceFileName)

	fmt.Printf(cli.Red+"🗑️  Menghapus file systemd %s...\n"+cli.Reset, serviceFilePath)
	err := os.Remove(serviceFilePath)
	if err != nil {
		return fmt.Errorf("Gagal menghapus file konfigurasi: %v", err)
	}

	fmt.Println(cli.Blue + "🔄 Merefresh systemd daemon..." + cli.Reset)
	if err := runSystemCommand("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("Gagal daemon-reload: %v", err)
	}

	fmt.Println(cli.Blue + "🧹 Membersihkan sisa status fail (reset-failed)..." + cli.Reset)
	_ = runSystemCommand("systemctl", "reset-failed")

	return nil
}

func ListManagedServices() {
	files, err := ioutil.ReadDir("/etc/systemd/system/")
	if err != nil {
		fmt.Printf(cli.Red+"❌ Gagal membaca direktori systemd: %v\n"+cli.Reset, err)
		return
	}

	fmt.Println(cli.Cyan + cli.Bold + "=====================================" + cli.Reset)
	fmt.Println(cli.Bold + "📋 Daftar Service Systemd (" + cli.Yellow + "Potensial milik Anda" + cli.Reset + cli.Bold + ")" + cli.Reset)
	fmt.Println(cli.Cyan + "=====================================" + cli.Reset)

	found := 0
	for _, file := range files {
		// Abaikan folder atau symlink (layanan bawaan sistem biasanya symlink dari /lib/)
		if file.Mode()&os.ModeSymlink != 0 || file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasSuffix(name, ".service") {
			// Read file content explicitly to see if it's managed by stackd, or just show regular files
			// In Linux, mostly user-added services reside here as physical files.
			fmt.Printf("   - "+cli.Green+"%s"+cli.Reset+"\n", strings.TrimSuffix(name, ".service"))
			found++
		}
	}

	if found == 0 {
		fmt.Println("   " + cli.Yellow + "Tidak ada service khusus/manual yang terdeteksi." + cli.Reset)
	}
	fmt.Println(cli.Cyan + "=====================================" + cli.Reset)
}

func runSystemCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
