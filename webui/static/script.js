function disableDarkMode() {
    document.body.classList.remove("dark-mode");
    localStorage.setItem("darkMode", null);
}

function enableDarkMode() {
    document.body.classList.add("dark-mode");
    localStorage.setItem("darkMode", "active");
}

function showBasedOnTheme() {
    let dark = localStorage.getItem("darkMode");
    dark == "active"
        ? enableDarkMode()
        : disableDarkMode();
}

function themeSwitch() {
    let dark = localStorage.getItem("darkMode");
    dark == "active"
        ? disableDarkMode()
        : enableDarkMode();
}

const serverHolder = document.getElementById("server-holder");
const socket = new WebSocket("ws://" + document.location.host + "/ws/");

const svg = `
<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
    <path strokeLinecap="round" strokeLinejoin="round" d="M5.25 14.25h13.5m-13.5 0a3 3 0 0 1-3-3m3 3a3 3 0 1 0 0 6h13.5a3 3 0 1 0 0-6m-16.5-3a3 3 0 0 1 3-3h13.5a3 3 0 0 1 3 3m-19.5 0a4.5 4.5 0 0 1 .9-2.7L5.737 5.1a3.375 3.375 0 0 1 2.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 0 1 .9 2.7m0 0a3 3 0 0 1-3 3m0 3h.008v.008h-.008v-.008Zm0-6h.008v.008h-.008v-.008Zm-3 6h.008v.008h-.008v-.008Zm0-6h.008v.008h-.008v-.008Z" />
</svg>`;

socket.onerror = function (e) {
    console.log("WebSocket error: ", e);
};

socket.onmessage = function (e) {
    let html = "";

    const serverData = JSON.parse(e.data)

    serverData.Active.forEach(server => {
        html += `
        <div class="server">
            ${svg}
            
            <p>${ server }</p>
            <p><span class="green-circle"}"></span> Online</p>
        </div>`;
    });

    serverData.Inactive.forEach(server => {
        html += `
        <div class="server">
            ${svg}
            
            <p>${ server }</p>
            <p><span class="red-circle"}"></span> Offline</p>
        </div>`;
    });

    serverHolder.innerHTML = html;
};

showBasedOnTheme();