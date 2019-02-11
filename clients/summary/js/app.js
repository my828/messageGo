var form = document.querySelector("form")
var input = document.querySelector('#input');

form.addEventListener("submit", evt => {
    evt.preventDefault();
    fetch("https://api.turtlemaster.me/v1/summary?url="+input.value)
    .then(handleResponse)
    .then(handleData)
    .catch(handleError)
})

var text = document.querySelector("p");
function handleData(data) {
    console.log(data);
    text.innerHTML = data;
}

function handleResponse(response) {
    if (response.ok) {
        return response.text();
    } else {
        return response.json()
            .then(err => {
                throw new Error(err.errorMessage);
            })
    }                                                   
}

function handleError(err) {
    alert(err.errorMessage)
}