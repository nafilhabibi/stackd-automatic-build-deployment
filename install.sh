#!/bin/bash

# Hentikan script jika terjadi error di salah satu perintah
set -e

echo -e "\033[1;36m=================================================\033[0m"
echo -e "\033[1;36m🚀 Memulai Instalasi stackd CLI Automator\033[0m"
echo -e "\033[1;36m=================================================\033[0m"

# Periksa apakah Golang sudah terinstall
if ! command -v go &> /dev/null
then
    echo -e "\033[0;31m❌ Error: Golang tidak ditemukan di sistem.\033[0m"
    echo "Karena stackd dibangun dari source code, pastikan Anda telah menginstal compiler Go (disarankan versi 1.20+)."
    echo "Panduan install: https://go.dev/doc/install"
    exit 1
fi

echo -e "\033[0;34m⏳ Menyiapkan modul dan mem-build binary stackd...\033[0m"
# Rapikan library
go mod tidy

# Eksekusi kompilasi
go build -o stackd cmd/stackd/main.go

echo -e "\033[0;32m✅ Build source code sukses!\033[0m"
echo ""
echo -e "\033[1;33m🔑 Meminta izin akses administrator (sudo) untuk memindahkan binary ke /usr/local/bin/...\033[0m"

# Pindahkan ke direktori global
sudo mv stackd /usr/local/bin/

echo -e "\033[1;36m=================================================\033[0m"
echo -e "\033[1;32m🎉 Instalasi stackd Selesai!\033[0m"
echo -e "Sekarang Anda dapat memanggil perintah \033[1mstackd\033[0m dari direktori mana saja."
echo -e "Coba ketik \033[1;33mstackd help\033[0m untuk memverifikasi instalasi."
echo -e "\033[1;36m=================================================\033[0m"
