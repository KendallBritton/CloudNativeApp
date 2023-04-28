function search_game() {
    let input = document.getElementById('searchbar').value;
    
    console.log(input);
    location.href = "./Resultpage.html?input="+input;
}

async function webpage_load(){
    let params = new URLSearchParams(document.location.search);
    let input = params.get("input");
    console.log(input);
    document.getElementById('test').innerHTML = "You typed in "+input;
    var data = {
        Name: input
    }
    console.log({
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    try {
        await fetch("/search", {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
            .then(response => {
                if (!response.ok) {
                    console.log("Not OK response!");
                    console.log(response)
                    return null;
                }
                return response.json()
            })
            .then(response => {
                if (response) {
                    console.log(response)
                    handleResponse(response)
                }
            })
            .catch(error => console.error(error));
    } catch (error) {
        console.log("Error:", error);
    }
}

var gameResponse;
function handleResponse(response) {
    gameResponse = response;
    // Game Title
    if (Object.keys(response.Name).length != 0){
        document.getElementById('name').innerHTML = response.Name;
    } else {
        document.getElementById('name').innerHTML = "Game Data Not Found In Database";
    }  
    console.log("Name = " + response.Name)

    // Game Story
    if (Object.keys(response.Storyline).length != 0){
        document.getElementById('storyline').innerHTML = response.Storyline;
    } else {
        document.getElementById('storyline').innerHTML = "Sorry! Storyline data was not found within the database!";
    }   
    console.log("Storyline = " + response.Storyline)


    // Boxart
    const boxartImage = new Image(350,350);
    boxartImage.src = response.Boxart;
    let boxart = document.getElementById("boxart");
    boxart.append(boxartImage);

    console.log("Cover url = " + response.Boxart)

    // Gameplay
    var tag = document.createElement('script');

    tag.src = "https://www.youtube.com/iframe_api";
    console.log(document.getElementsByTagName('script'))
    var firstScriptTag = document.getElementsByTagName('script')[0];
    firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

}

var player;
const gameplayIDs = ['gameplay-1','gameplay-2','gameplay-3','gameplay-4','gameplay-5','gameplay-6','gameplay-7','gameplay-8','gameplay-9','gameplay-10'];
function onYouTubeIframeAPIReady() {
    for (let i = 0; i < gameplayIDs.length; i++){
        console.log(gameResponse.VideoId);
        player = new YT.Player(gameplayIDs[i], {
            height: '350',
            width: '365',
            videoId: gameResponse.VideoId[i]
        });
    }
}

function isEmptyObjectStoryline(response) {
  return response.Storyline && Object.keys(response.Storyline).length === 0 && response.Storyline.constructor === Object;
}