/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** ServiceData
*/

import Netflix from "./../assets/netflix.png"
import Spotify from "./../assets/spotify.png"
import Instagram from "./../assets/instagram.webp"
import X from "./../assets/x.webp"
import Nasa from "./../assets/nasa.webp"

const services = [
    {
        name: "Spotify",
        logo: Spotify,
        link: "https://www.spotify.com",
        color: "bg-spotify-100 hover:bg-spotify-200"
    },
    {
        name: "Netflix",
        logo: Netflix,
        link: "https://www.netflix.com",
        color: "bg-red-600 hover:bg-red-500"
    },
    {
        name: "Instagram",
        logo: Instagram,
        link: "https://www.instagram.com",
        color: "bg-instagram-100 hover:bg-instagram-200"
    },
    {
        name: "Twitter",
        logo: X,
        link: "https://www.x.com",
        color: "bg-gray-900 hover:bg-gray-800"
    },
    {
        name: "Nasa",
        logo: Nasa,
        link: "https://www.nasa.gov",
        color: "bg-nasa-100 hover:bg-nasa-200"
    }
]

export default services