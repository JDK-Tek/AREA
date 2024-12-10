/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** fetchData
*/

export default async function fetchData(url, request = {method: "GET"}) {

    try {
        const res = await fetch(url, request);
        if (!res.ok) {
            throw res.statusText;
        }
    
        const json = await res.json();
        return { success: true, data: json };
    } catch (err) {
        return { success: false, error: err };
    }
}
