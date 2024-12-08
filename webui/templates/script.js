const port = document.getElementById("helper").getAttribute("port");
const serverHolder = document.getElementById("server-holder");

let socket = new WebSocket(`ws://localhost:${port}/ws/`);

socket.onmessage = function (e) {
    e.forEach((server) => console.log(`${server.URL} ${server.IsAlive}`)) 
};