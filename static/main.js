const refreshGames = (event) => {
    event.preventDefault();
    const url = "/refresh-games";
    console.log(event.target.elements.playerId.value);
    fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            playerId: event.target.elements.playerId.value
        })
    }).then(res => {
        return res.text()
    }).then(body => {
        const refreshResponseP = document.querySelector('#refresh-response');
        refreshResponseP.innerHTML = body;
        refreshResponseP.style.display = "block";
    });
}

const form = document.getElementById('refresh-games-form');
form.addEventListener('submit', refreshGames);

const getPlayers = () => {
    const url = "/get-players";
    fetch(url, {
        method: 'GET'
    }).then(res => {
        return res.text()
    }).then(body => {
        console.log(body);
    });
}