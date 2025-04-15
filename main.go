package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

/*
GOOS=windows GOARCH=amd64 go build -o app.exe main.go
GOOS=linux GOARCH=amd64 go build -o app main.go
GOOS=darwin GOARCH=amd64 go build -o app-darwin main.go
*/

const htmlPage = `<!DOCTYPE html><html lang="en">

<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>QR Code Generator</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" rel="stylesheet" />
  <style>
    @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&display=swap');

    .glass-effect {
      backdrop-filter: blur(10px);
      background: rgba(255, 255, 255, 0.7);
    }

    .animate-float {
      animation: float 6s ease-in-out infinite;
    }

    @keyframes float {
      0% {
        transform: translateY(0px);
      }

      50% {
        transform: translateY(-10px);
      }

      100% {
        transform: translateY(0px);
      }
    }

    .card-shadow {
      box-shadow: 0 10px 30px rgba(139, 92, 246, 0.1);
    }

    .btn-shine {
      position: relative;
      overflow: hidden;
    }

    .btn-shine::after {
      content: '';
      position: absolute;
      top: -50%;
      left: -50%;
      width: 200%;
      height: 200%;
      background: linear-gradient(to right, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 0.3) 50%, rgba(255, 255, 255, 0) 100%);
      transform: rotate(30deg);
      transition: 0.6s;
    }

    .btn-shine:hover::after {
      left: 100%;
    }

    .size-option:checked+label {
      background: linear-gradient(to right, #8b5cf6, #d946ef);
      color: white;
      border-color: transparent;
    }
  </style>
</head>

<body
  class="bg-gradient-to-br from-indigo-100 via-purple-100 to-pink-100 min-h-screen flex items-center justify-center p-4 font-[Poppins]">
  <!-- Background Decorations -->
  <div class="fixed inset-0 overflow-hidden pointer-events-none z-0">
    <div
      class="absolute top-20 left-20 w-64 h-64 rounded-full bg-gradient-to-r from-pink-200 to-pink-300 opacity-30 blur-3xl animate-float">
    </div>
    <div
      class="absolute bottom-20 right-20 w-80 h-80 rounded-full bg-gradient-to-r from-violet-200 to-indigo-300 opacity-30 blur-3xl animate-float"
      style="animation-delay: 2s;"></div>
    <div
      class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 rounded-full bg-gradient-to-r from-blue-200 to-cyan-300 opacity-20 blur-3xl animate-float"
      style="animation-delay: 4s;"></div>
  </div>

  <!-- Main Card -->
  <div
    class="max-w-md w-full bg-white/90 backdrop-blur-lg rounded-3xl shadow-xl overflow-hidden card-shadow relative z-10 border border-white/50">
    <!-- Header -->
    <div class="bg-gradient-to-r from-violet-600 to-fuchsia-600 p-6 text-white relative overflow-hidden">
      <div class="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-bl-full"></div>
      <div class="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-tr-full"></div>

      <div class="flex items-center justify-between relative z-10">
        <h1 class="text-2xl font-bold flex items-center">
          <i class="fas fa-qrcode mr-3 text-3xl"></i> QR Code Generator
        </h1>
        <button id="authorInfoBtn"
          class="h-10 w-10 bg-white/20 rounded-full flex items-center justify-center hover:bg-white/30 transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-white/50"
          aria-label="About Author">
          <i class="fas fa-info text-white"></i>
        </button>
      </div>
      <p class="mt-2 opacity-90 font-light">Create, download and share QR codes instantly</p>
    </div>

    <!-- Form -->
    <div class="p-6 space-y-5">
      <div class="group transition-all duration-300 hover:translate-x-1">
        <label for="filename"
          class="block text-sm font-medium text-gray-700 mb-1 group-hover:text-violet-600 transition-colors">
          <i class="fas fa-file-alt mr-2"></i>File name
        </label>
        <input type="text" id="filename" placeholder="my-qr-code"
          class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-violet-500 focus:border-transparent outline-none transition-all duration-300 hover:border-violet-300" />
      </div>

      <div class="group transition-all duration-300 hover:translate-x-1">
        <label for="qrdata"
          class="block text-sm font-medium text-gray-700 mb-1 group-hover:text-violet-600 transition-colors">
          <i class="fas fa-link mr-2"></i>QR Content
        </label>
        <input type="text" id="qrdata" placeholder="https://example.com"
          class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-violet-500 focus:border-transparent outline-none transition-all duration-300 hover:border-violet-300" />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="group transition-all duration-300 hover:translate-x-1">
          <label for="qrcolor"
            class="block text-sm font-medium text-gray-700 mb-1 group-hover:text-violet-600 transition-colors">
            <i class="fas fa-palette mr-2"></i>QR Color
          </label>
          <div class="relative">
            <input type="color" id="qrcolor" value="#000000"
              class="w-full h-12 px-3 py-2 border border-gray-200 rounded-xl cursor-pointer transition-all duration-300 hover:border-violet-300" />
            <div class="absolute inset-y-0 right-3 flex items-center pointer-events-none">
              <span class="text-xs font-medium text-gray-500" id="colorHex">#000000</span>
            </div>
          </div>
        </div>

        <div class="group">
          <label for="qrsize"
            class="block text-sm font-medium text-gray-700 mb-1 group-focus-within:text-violet-600 transition-colors flex items-center">
            <i class="fas fa-expand-arrows-alt mr-2 text-violet-500"></i>QR Size
          </label>
          <select id="qrsize"
            class="w-full h-12 px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-violet-500 focus:border-transparent outline-none input-transition appearance-none bg-white">
            <option value="200">Small (200px)</option>
            <option value="300" selected>Medium (300px)</option>
            <option value="400">Large (400px)</option>
            <option value="500">X-Large (500px)</option>
          </select>
        </div>
      </div>

      <button id="generateBtn"
        class="w-full bg-gradient-to-r from-violet-600 to-fuchsia-600 text-white font-medium py-3 px-4 rounded-xl hover:from-violet-700 hover:to-fuchsia-700 transition transform hover:scale-[1.02] active:scale-[0.98] shadow-lg hover:shadow-xl btn-shine">
        <i class="fas fa-magic mr-2"></i>Generate QR Code
      </button>
    </div>

    <!-- Result -->
    <div id="output" class="p-6 border-t border-gray-100 hidden">
      <div class="flex flex-col items-center">
        <div
          class="relative mb-5 w-full max-w-xs mx-auto bg-gradient-to-br from-gray-50 to-white p-4 rounded-xl shadow-inner">
          <div id="loadingOverlay"
            class="absolute inset-0 bg-white/90 flex items-center justify-center rounded-xl hidden z-10">
            <div class="flex flex-col items-center">
              <div class="animate-spin rounded-full h-12 w-12 border-4 border-violet-200 border-t-violet-600"></div>
              <p class="mt-3 text-sm text-violet-600 font-medium">Generating...</p>
            </div>
          </div>
          <img id="qrImage" class="w-full h-auto rounded-lg shadow-sm transition-all duration-500 hover:scale-105"
            alt="Generated QR Code" />
        </div>

        <div class="grid grid-cols-2 gap-4 w-full">
          <a id="downloadBtn"
            class="bg-gradient-to-r from-sky-500 to-blue-500 hover:from-sky-600 hover:to-blue-600 text-white font-medium py-3 px-4 rounded-xl text-center flex items-center justify-center transform hover:scale-[1.02] active:scale-[0.98] shadow-lg hover:shadow-xl transition-all duration-300 btn-shine"
            download>
            <i class="fas fa-download mr-2"></i>Download
          </a>
          <button id="shareBtn"
            class="bg-gradient-to-r from-violet-500 to-fuchsia-500 hover:from-violet-600 hover:to-fuchsia-600 text-white font-medium py-3 px-4 rounded-xl flex items-center justify-center transform hover:scale-[1.02] active:scale-[0.98] shadow-lg hover:shadow-xl transition-all duration-300 btn-shine">
            <i class="fas fa-share-alt mr-2"></i>Share
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Toast Notification -->
  <div id="toast"
    class="fixed bottom-4 right-4 bg-gray-800/90 backdrop-blur-sm text-white px-5 py-3 rounded-xl shadow-lg transform translate-y-20 opacity-0 transition-all duration-300 flex items-center max-w-xs z-50">
    <i id="toastIcon" class="mr-2"></i>
    <span id="toastMessage"></span>
  </div>

  <!-- Author Info Modal -->
  <div id="authorModal"
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 opacity-0 pointer-events-none transition-opacity duration-300 z-50">
    <div
      class="bg-white/95 backdrop-blur-md rounded-3xl shadow-2xl max-w-md w-full transform scale-95 transition-transform duration-300 border border-white/50"
      id="modalContent">
      <div class="relative">
        <!-- Modal Header with Gradient -->
        <div
          class="bg-gradient-to-r from-violet-600 to-fuchsia-600 rounded-t-3xl p-6 text-white relative overflow-hidden">
          <div class="absolute top-0 right-0 w-24 h-24 bg-white/10 rounded-bl-full"></div>
          <div class="flex items-center justify-between relative z-10">
            <h2 class="text-xl font-bold flex items-center">
              <i class="fas fa-user-circle text-2xl mr-3"></i> About the Developer
            </h2>
            <button id="closeModalBtn"
              class="h-8 w-8 bg-white/20 rounded-full flex items-center justify-center hover:bg-white/30 transition-all focus:outline-none transform hover:scale-105"
              aria-label="Close dialog">
              <i class="fas fa-times text-white"></i>
            </button>
          </div>
        </div>

        <!-- Author Info Content -->
        <div class="p-6">
          <div class="flex flex-col md:flex-row items-center mb-6">
            <div
              class="h-24 w-24 rounded-full bg-gradient-to-r from-violet-500 to-fuchsia-500 flex items-center justify-center text-white text-3xl shadow-lg transform hover:scale-105 transition-transform duration-300">
              <i class="fas fa-user"></i>
            </div>
            <div class="mt-4 md:mt-0 md:ml-4 text-center md:text-left">
              <h3 class="text-xl font-bold text-gray-800">Saravanan Selvam</h3>
              <p class="text-gray-600">Full Stack Developer</p>
              <div class="flex justify-center md:justify-start space-x-3 mt-2">
                <a href="https://github.com/saravanan611" class="text-gray-500 hover:text-violet-600 transition-colors">
                  <i class="fab fa-github"></i>
                </a>
                <a href="linkedin.com/in/saravanan-s-142ab11bb/" class="text-gray-500 hover:text-blue-600 transition-colors">
                  <i class="fab fa-linkedin"></i>
                </a>
              </div>
            </div>
          </div>

          <div class="space-y-4">
            <div
              class="flex items-center p-4 bg-gray-50/80 backdrop-blur-sm rounded-xl hover:bg-gray-100/80 transition-colors duration-300 transform hover:translate-x-1">
              <div
                class="h-12 w-12 rounded-full bg-gradient-to-r from-violet-500/20 to-fuchsia-500/20 flex items-center justify-center text-violet-600">
                <i class="fas fa-envelope"></i>
              </div>
              <div class="ml-4">
                <p class="text-sm text-gray-500">Email</p>
                <p class="text-gray-800 font-medium">saravanan9944gs@gmail.com</p>
              </div>
            </div>

            <div
              class="flex items-center p-4 bg-gray-50/80 backdrop-blur-sm rounded-xl hover:bg-gray-100/80 transition-colors duration-300 transform hover:translate-x-1">
              <div
                class="h-12 w-12 rounded-full bg-gradient-to-r from-violet-500/20 to-fuchsia-500/20 flex items-center justify-center text-violet-600">
                <i class="fas fa-phone"></i>
              </div>
              <div class="ml-4">
                <p class="text-sm text-gray-500">Phone</p>
                <p class="text-gray-800 font-medium">+91 9080921942</p>
              </div>
            </div>
          </div>

          <div class="mt-6 text-center">
            <p class="text-sm text-gray-500">© 2025 QR Code Generator. All rights reserved.</p>
          </div>
        </div>
      </div>
    </div>
  </div>

  <script>
    const elements = {
      filename: document.getElementById("filename"),
      qrdata: document.getElementById("qrdata"),
      qrcolor: document.getElementById("qrcolor"),
      qrsize: document.getElementById("qrsize"),
      colorHex: document.getElementById("colorHex"),
      generateBtn: document.getElementById("generateBtn"),
      downloadBtn: document.getElementById("downloadBtn"),
      shareBtn: document.getElementById("shareBtn"),
      output: document.getElementById("output"),
      qrImage: document.getElementById("qrImage"),
      loadingOverlay: document.getElementById("loadingOverlay"),
      toast: document.getElementById("toast"),
      toastIcon: document.getElementById("toastIcon"),
      toastMessage: document.getElementById("toastMessage"),
      authorInfoBtn: document.getElementById("authorInfoBtn"),
      authorModal: document.getElementById("authorModal"),
      closeModalBtn: document.getElementById("closeModalBtn")
    };

    const toastTypes = {
      success: { icon: "fas fa-check-circle text-green-400" },
      error: { icon: "fas fa-exclamation-circle text-red-400" },
      info: { icon: "fas fa-info-circle text-blue-400" }
    };

    const showToast = (msg, type = "info") => {
      elements.toastMessage.textContent = msg;
      elements.toastIcon.className = toastTypes[type].icon + " mr-2";
      elements.toast.classList.remove("translate-y-20", "opacity-0");
      elements.toast.classList.add("translate-y-0", "opacity-100");
      setTimeout(() => {
        elements.toast.classList.add("translate-y-20", "opacity-0");
        elements.toast.classList.remove("translate-y-0", "opacity-100");
      }, 3000);
    };

    let qrBlobURL = null;

    function toggleLoading(isLoading) {
      elements.loadingOverlay.classList.toggle("hidden", !isLoading);
      elements.generateBtn.disabled = isLoading;
      elements.generateBtn.classList.toggle("opacity-70", isLoading);
    }

    function getSelectedSize() {
      if (elements.sizeSmall.checked) return elements.sizeSmall.value;
      if (elements.sizeMedium.checked) return elements.sizeMedium.value;
      if (elements.sizeLarge.checked) return elements.sizeLarge.value;
      if (elements.sizeXlarge.checked) return elements.sizeXlarge.value;
      return "250"; // Default to medium if somehow none are selected
    }

    async function generateQR() {
      const filename = elements.filename.value.trim(),
      qrdata = elements.qrdata.value.trim(),
      qrcolor = elements.qrcolor.value,
      qrsize = elements.qrsize.value;

      if (!filename) return showToast("Please enter a filename", "error");
      if (!qrdata) return showToast("Please enter QR code data", "error");

      toggleLoading(true);

      try {
        const response = await fetch("/qr/Generate", {
          method: "GET",
          headers: {
            filename: filename,
            qrdata: qrdata,
            qrcolor: qrcolor,
            qrsize: qrsize
          }
        });

        if (!response.ok) throw new Error("Failed to generate QR code");

        if (qrBlobURL) URL.revokeObjectURL(qrBlobURL);
        const blob = await response.blob();
        qrBlobURL = URL.createObjectURL(blob);

        elements.qrImage.src = qrBlobURL;
        elements.output.classList.remove("hidden");
        elements.downloadBtn.href = qrBlobURL;
        elements.downloadBtn.download = filename + ".png";

        elements.shareBtn.onclick = () => {
          if (navigator.share) {
            navigator.share({
              title: filename,
              text: "Here's your QR code!",
              files: [new File([blob], filename + ".png", { type: blob.type })]
            }).catch(e => showToast("Share failed: " + e, "error"));
          } else {
            showToast("Sharing not supported on this browser", "error");
          }
        };

        showToast("QR Code generated!", "success");
      } catch (err) {
        console.error(err);
        showToast("Error generating QR code", "error");
      } finally {
        toggleLoading(false);
      }
    }

    // Update color hex display
    elements.qrcolor.addEventListener("input", () => {
      elements.colorHex.textContent = elements.qrcolor.value;
    });

    elements.generateBtn.addEventListener("click", generateQR);
    elements.authorInfoBtn.addEventListener("click", () => {
      elements.authorModal.classList.remove("opacity-0", "pointer-events-none");
      elements.authorModal.classList.add("opacity-100", "pointer-events-auto");
    });
    elements.closeModalBtn.addEventListener("click", () => {
      elements.authorModal.classList.add("opacity-0", "pointer-events-none");
      elements.authorModal.classList.remove("opacity-100", "pointer-events-auto");
    });
  </script>

</body>

</html>`

func main() {

	lPort := os.Getenv("Qrport")
	if lPort == "" {
		lPort = "1234"
	}

	lUrl := "http://localhost:" + lPort

	if lErr := openBrowser(lUrl + "/qr/open/"); lErr != nil {
		Notification("openbrowser error", lErr.Error(), "")
		os.Exit(1)
	}

	http.HandleFunc("/qr/open/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			fmt.Fprint(w, "optional call success")
			return
		}

		fmt.Fprint(w, htmlPage)
	})
	http.HandleFunc("/qr/Generate", CreateQr)

	time.AfterFunc(time.Duration(time.Minute*60), func() {
		Notification("hower pass", "qr app life time is 1 hr < re-open the app evert 1 hour >", "")
		os.Exit(0)
	})

	if lErr := http.ListenAndServe(":"+lPort, nil); lErr != nil {
		Notification("port run error", lErr.Error(), "")
		os.Exit(1)
	}

}

func Notification(lTitle, lMsg, lUrl string) (lErr error) {
	log.Println("Error :", lMsg)
	if lUrl == "" {
		lUrl = " contact < saravanan9944gs@gmail.com > "
	}
	var lCmd *exec.Cmd
	lCurrentTime := time.Now().Format("2006-01-02 15:04:05")
	switch runtime.GOOS {
	case "windows":
		lCmd = exec.Command("powershell", "-Command",
			"New-BurntToastNotification -Text '"+lTitle+"', '"+lMsg+" "+lUrl+" - Time: "+lCurrentTime+"'")
	case "darwin":
		lCmd = exec.Command("osascript", "-e",
			`display notification "`+lMsg+` `+lUrl+` - Time: `+lCurrentTime+`" with title "`+lTitle+`" sound name "default"`)
	case "linux":
		lCmd = exec.Command("notify-send", lTitle, lMsg+"\n"+lUrl+"\n - Time: "+lCurrentTime)
	default:
		log.Fatalf("unsupported platform: %s", runtime.GOOS)
	}

	return lCmd.Run()
}

func CreateQr(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "false")
	(w).Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "qrdata,filename,qrcolor,qrsize,Accept,user,INDICATOR,PROGRAMID,SEARCHDATE,isActive,schedulerId,repoid,Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	if r.Method == http.MethodOptions {
		fmt.Fprint(w, "optional call success")
		return
	}

	var lErr error

	defer func() {
		if lErr != nil {
			fmt.Fprint(w, "error in responce")
			Notification("QrCreating Error", lErr.Error(), "")
		}
	}()

	lQrContent := r.Header.Get("qrdata")
	if lQrContent == "" {
		lErr = fmt.Errorf("qr msg value is missing")
		return
	}

	lSize, lErr := strconv.Atoi(r.Header.Get("qrsize"))
	if lErr != nil {
		return
	}
	lQrByte, lErr := QrByte(lQrContent, r.Header.Get("qrcolor"), lSize)
	if lErr != nil {
		return
	}

	w.Header().Set("filename", "simpleQR.png")
	w.Header().Set("content-type", "image/png")
	_, lErr = w.Write(lQrByte)
	if lErr != nil {
		return
	}

}

func hexToRGBA(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.RGBA{}, errors.New("invalid hex color format")
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

// func QrByte(lUrl, lColor string, lSize int) (lQrByte []byte, lErr error) {
// 	return qrcode.Encode(lUrl, qrcode.Medium, lSize)
// }

func QrByte(lUrl, lColor string, lSize int) ([]byte, error) {
	qr, err := qrcode.New(lUrl, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	// Convert hex color to RGBA
	rgba, err := hexToRGBA(lColor)
	if err != nil {
		return nil, err
	}

	qr.ForegroundColor = rgba
	qr.BackgroundColor = color.White // or customize

	// Encode to PNG and return bytes
	var buf bytes.Buffer
	err = qr.Write(lSize, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		// cmd = "rundll32"
		// args = []string{"url.dll,FileProtocolHandler", url}
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}
