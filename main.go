package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/skip2/go-qrcode"
)

/*
GOOS=windows GOARCH=amd64 go build -o app.exe main.go
GOOS=linux GOARCH=amd64 go build -o app main.go
GOOS=darwin GOARCH=amd64 go build -o app-darwin main.go
*/

const htmlPage = `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"/><meta name="viewport" content="width=device-width, initial-scale=1.0"/><title>QR Code Generator</title><script src="https://cdn.tailwindcss.com"></script><link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" rel="stylesheet"/></head><body class="bg-gradient-to-br from-sky-100 to-violet-100 min-h-screen flex items-center justify-center p-4 font-sans"><div class="max-w-md w-full bg-white rounded-2xl shadow-xl overflow-hidden"><div class="bg-gradient-to-r from-violet-500 to-fuchsia-500 p-6 text-white"><div class="flex items-center justify-between"><h1 class="text-2xl font-bold flex items-center"><i class="fas fa-qrcode mr-3 text-3xl"></i> QR Code Generator</h1><button id="authorInfoBtn" class="h-10 w-10 bg-white bg-opacity-20 rounded-full flex items-center justify-center hover:bg-opacity-30" aria-label="About Author"><i class="fas fa-info text-white"></i></button></div><p class="mt-2 opacity-90">Create, download and share QR codes instantly</p></div><div class="p-6 space-y-5"><div><label for="filename" class="block text-sm font-medium text-gray-700 mb-1"><i class="fas fa-file-alt mr-2"></i>File name</label><input type="text" id="filename" placeholder="my-qr-code" class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-violet-500 outline-none"/></div><div><label for="qrdata" class="block text-sm font-medium text-gray-700 mb-1"><i class="fas fa-link mr-2"></i>QR Content</label><input type="text" id="qrdata" placeholder="https://example.com" class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:ring-2 focus:ring-violet-500 outline-none"/></div><button id="generateBtn" class="w-full bg-gradient-to-r from-violet-500 to-fuchsia-500 text-white font-medium py-3 px-4 rounded-xl hover:from-violet-600 hover:to-fuchsia-600 transition transform hover:scale-105"><i class="fas fa-magic mr-2"></i>Generate QR Code</button></div><div id="output" class="p-6 border-t border-gray-100 hidden"><div class="flex flex-col items-center"><div class="relative mb-5 w-full max-w-xs mx-auto bg-gray-50 p-3 rounded-xl"><div id="loadingOverlay" class="absolute inset-0 bg-white bg-opacity-90 flex items-center justify-center rounded-xl hidden z-10"><div class="animate-spin rounded-full h-12 w-12 border-4 border-violet-200 border-t-violet-600"></div></div><img id="qrImage" class="w-full h-auto rounded-lg shadow-sm" alt="Generated QR Code"/></div><div class="grid grid-cols-2 gap-4 w-full"><a id="downloadBtn" class="bg-gradient-to-r from-sky-500 to-blue-500 hover:from-sky-600 hover:to-blue-600 text-white font-medium py-3 px-4 rounded-xl text-center flex items-center justify-center transform hover:scale-105" download><i class="fas fa-download mr-2"></i>Download</a><button id="shareBtn" class="bg-gradient-to-r from-violet-500 to-fuchsia-500 hover:from-violet-600 hover:to-fuchsia-600 text-white font-medium py-3 px-4 rounded-xl flex items-center justify-center transform hover:scale-105"><i class="fas fa-share-alt mr-2"></i>Share</button></div></div></div></div><div id="toast" class="fixed bottom-4 right-4 bg-gray-800 text-white px-5 py-3 rounded-xl shadow-lg transform translate-y-20 opacity-0 transition-all duration-300 flex items-center max-w-xs z-50"><i id="toastIcon" class="mr-2"></i><span id="toastMessage"></span></div><div id="authorModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 opacity-0 pointer-events-none transition-opacity duration-300 z-50"><div class="bg-white rounded-2xl shadow-2xl max-w-md w-full transform scale-95 transition-transform duration-300" id="modalContent"><div class="bg-gradient-to-r from-violet-500 to-fuchsia-500 rounded-t-2xl p-6 text-white"><div class="flex items-center justify-between"><h2 class="text-xl font-bold flex items-center"><i class="fas fa-user-circle text-2xl mr-3"></i> About the Developer</h2><button id="closeModalBtn" class="h-8 w-8 bg-white bg-opacity-20 rounded-full flex items-center justify-center hover:bg-opacity-30" aria-label="Close dialog"><i class="fas fa-times text-white"></i></button></div></div><div class="p-6"><div class="flex items-center mb-6"><div class="h-20 w-20 rounded-full bg-gradient-to-r from-violet-500 to-fuchsia-500 flex items-center justify-center text-white text-3xl"><i class="fas fa-user"></i></div><div class="ml-4"><h3 class="text-xl font-bold text-gray-800">Saravanan Selvam</h3><p class="text-gray-600">Full Stack Developer</p></div></div><div class="space-y-4"><div class="flex items-center p-3 bg-gray-50 rounded-xl"><div class="h-10 w-10 rounded-full bg-violet-100 flex items-center justify-center text-violet-600"><i class="fas fa-envelope"></i></div><div class="ml-3"><p class="text-sm text-gray-500">Email</p><p class="text-gray-800">saravanan9944gs@gmail.com</p></div></div><div class="flex items-center p-3 bg-gray-50 rounded-xl"><div class="h-10 w-10 rounded-full bg-violet-100 flex items-center justify-center text-violet-600"><i class="fas fa-phone"></i></div><div class="ml-3"><p class="text-sm text-gray-500">Phone</p><p class="text-gray-800">+91 9080921942</p></div></div></div></div></div></div><script>const elements={filename:document.getElementById("filename"),qrdata:document.getElementById("qrdata"),generateBtn:document.getElementById("generateBtn"),downloadBtn:document.getElementById("downloadBtn"),shareBtn:document.getElementById("shareBtn"),output:document.getElementById("output"),qrImage:document.getElementById("qrImage"),loadingOverlay:document.getElementById("loadingOverlay"),toast:document.getElementById("toast"),toastIcon:document.getElementById("toastIcon"),toastMessage:document.getElementById("toastMessage"),authorInfoBtn:document.getElementById("authorInfoBtn"),authorModal:document.getElementById("authorModal"),closeModalBtn:document.getElementById("closeModalBtn")},toastTypes={success:{icon:"fas fa-check-circle text-green-400"},error:{icon:"fas fa-exclamation-circle text-red-400"},info:{icon:"fas fa-info-circle text-blue-400"}},showToast=(e,t="info")=>{elements.toastMessage.textContent=e,elements.toastIcon.className=toastTypes[t].icon+" mr-2",elements.toast.classList.remove("translate-y-20","opacity-0"),elements.toast.classList.add("translate-y-0","opacity-100"),setTimeout(()=>{elements.toast.classList.add("translate-y-20","opacity-0"),elements.toast.classList.remove("translate-y-0","opacity-100")},3e3)};let qrBlobURL=null;function toggleLoading(e){elements.loadingOverlay.classList.toggle("hidden",!e),elements.generateBtn.disabled=e,elements.generateBtn.classList.toggle("opacity-70",e)}async function generateQR(){const e=elements.filename.value.trim(),t=elements.qrdata.value.trim();if(!e)return void showToast("Please enter a filename","error");if(!t)return void showToast("Please enter QR code data","error");toggleLoading(!0);try{const r=await fetch("/qr/Generate",{method:"GET",headers:{filename:e,qrdata:t}});if(!r.ok)throw new Error("Failed to generate QR code");qrBlobURL&&URL.revokeObjectURL(qrBlobURL);const o=await r.blob();qrBlobURL=URL.createObjectURL(o),elements.qrImage.src=qrBlobURL,elements.output.classList.remove("hidden"),elements.downloadBtn.href=qrBlobURL,elements.downloadBtn.download=e+".png",elements.shareBtn.onclick=()=>{navigator.share?navigator.share({title:e,text:"Here’s your QR code!",files:[new File([o],e+".png",{type:o.type})]}).catch((e=>showToast("Share failed: "+e,"error"))):showToast("Sharing not supported on this browser","error")},showToast("QR Code generated!","success")}catch(e){console.error(e),showToast("Error generating QR code","error")}finally{toggleLoading(!1)}}elements.generateBtn.addEventListener("click",generateQR),elements.authorInfoBtn.addEventListener("click",()=>{elements.authorModal.classList.remove("opacity-0","pointer-events-none"),elements.authorModal.classList.add("opacity-100","pointer-events-auto")}),elements.closeModalBtn.addEventListener("click",()=>{elements.authorModal.classList.add("opacity-0","pointer-events-none"),elements.authorModal.classList.remove("opacity-100","pointer-events-auto")});</script></body></html>`

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
	(w).Header().Set("Access-Control-Allow-Headers", "qrdata,filename,Accept,user,INDICATOR,PROGRAMID,SEARCHDATE,isActive,schedulerId,repoid,Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

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

	lQrByte, lErr := QrByte(lQrContent)
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

func QrByte(lUrl string) (lQrByte []byte, lErr error) {
	return qrcode.Encode(lUrl, qrcode.Medium, 256)
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
