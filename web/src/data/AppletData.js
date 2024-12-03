/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** AppletData
*/

import Netflix from "./../assets/netflix.png"
import Spotify from "./../assets/spotify.png"

import WeatherUnderground from "./../assets/weather-underground.webp"
import Notification from "./../assets/notification.webp"

import Instagram from "./../assets/instagram.webp"
import X from "./../assets/x.webp"

const AppletData = [
    {
        title: "Create playlist of your favorite series in one click",
        users : 132124,
        color: "bg-spotify-100 hover:bg-spotify-200",
        link: "https://spotify.com",
        services : [
            {
                name: "Spotify",
                logo: Spotify
            },
            {
                name: "Netflix",
                logo: Netflix
            }
        ]
    },
    {
        title: "Get the weather forecast every day at 7:00 AM",
        users : 88432,
        color: "bg-weatherunderground-100 hover:bg-weatherunderground-200",
        link: "https://www.wunderground.com/",
        services : [
            {
                name: "Weather Underground",
                logo: WeatherUnderground
            },
            {
                name: "Notification",
                logo: Notification
            }
        ]
    },
    {
        title: "Tweet your Instagrams as native photos on Twitter",
        users : 603723,
        color: "bg-instagram-100 hover:bg-instagram-200",
        link: "https://instagram.com",
        services : [
            {
                name: "Instagram",
                logo: Instagram
            },
            {
                name: "X",
                logo: X
            }
        ]
    },
];

export default AppletData;
