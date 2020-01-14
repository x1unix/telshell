#!/bin/sh
PKG_URL="github.com/x1unix/telshell"
URL_DOWNLOAD_PREFIX="https://${PKG_URL}/releases/latest/download"
ISSUE_URL="https://${PKG_URL}/issues"

RED="\033[0;31m"
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

warn() {
    printf "${YELLOW}${1}${NC}\n"
}

panic() {
    printf "${RED}ERROR: ${1}${NC}\n" >&2
    printf "${RED}\nIf you feel that this is an installer issue, you may submit an issue on ${ISSUE_URL}\nInstallation failed. ${NC}\n"
    exit 1
}

get_bin_name() {
    os=$(uname -s | awk '{print tolower($0)}')
    case $os in
    cygwin*|mingw32*|msys*|mingw*)
      os=windows
      file_suffix=.exe
      ;;
    *)
      ;;
    esac

    arc=$(get_arch)
    echo "telshell_${os}-${arc}${file_suffix}"
}

get_arch() {
    a=$(uname -m)
    case ${a} in
    "x86_64" | "amd64" )
      echo "amd64"
      ;;
    "i386" | "i486" | "i586" | "i686")
      echo "i386"
      ;;
    *)
      echo $a
      ;;
    esac
}


main() {
    file_name=$(get_bin_name)
    download_dir="${HOME}/bin"
    mkdir -p ${download_dir}

    dest_file="${download_dir}/telshell"
    download_url=${URL_DOWNLOAD_PREFIX}/${file_name}
    echo "-> Downloading '${download_url}'..."

    http_status=$(curl --fail --write-out "%{http_code}" -L --show-error --progress -o ${dest_file} ${download_url})
    case ${http_status} in
    "200")
      chmod +x ${dest_file}
      echo "-> Successfully installed to '${dest_file}'"
      printf "${GREEN}Done!${NC}\n"
      exit 0
      ;;
    "404")
      sys_label="$(uname -s) $(uname -m)"
      panic "No prebuilt binaries available for ${sys_label}, try to check out release for your platform at https://${PKG_URL}/releases"
      ;;
    *)
      panic "Installation failed, failed to download binary (HTTP error ${http_status})"
      ;;
    esac
}

main