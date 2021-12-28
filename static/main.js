const refreshGames = () => {
    const url = "/refresh-games";
    fetch(url, {
        method: 'POST'
    }).then(res => {
        return res.text()
    }).then(body => {
        const refreshResponseP = document.querySelector('#refresh-response');
        refreshResponseP.innerHTML = body;
        refreshResponseP.style.display = "block";
    });
}