/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** login
*/


var startTime = performance.now();
window.onload = function() {
    var endTime = performance.now();
    var loadTime = (endTime - startTime);
    document.getElementById('load-time').textContent = 'Temps de chargement: ' + loadTime.toFixed(2) + ' ms';
};
     

async function validateForm(event) {
    event.preventDefault();
    const email = document.getElementById('txtbox-email').value;
    const password = document.getElementById('txtbox-password').value;

    if (!email || !password) {
        document.getElementById('response').innerText = 'Please fill all fields.';
        return;
    }

    try {
        // const response = await fetch('http://localhost:1234/api/login', {
        //     method: 'POST',
        //     headers: {
        //         'Content-Type': 'application/json'
        //     },
        //     body: JSON.stringify({ email, password })
        // });

        // if (!response.ok) {
        //     throw new Error(`Network error: ${response.status} ${response.statusText}`);
        // }

        // const data = await response.json();
        // console.log('Response:', data);

        localStorage.setItem('email', email);
        // localStorage.setItem('token', data.token);
        window.location.href = 'content.html';
    } catch (error) {
        console.error('Erreur:', error.message);
        document.getElementById('response').innerText = `Error during API request: ${error.message}`;
    }
}
