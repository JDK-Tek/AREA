/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Explore
*/

import HeaderBar from "../../components/Header/HeaderBar";
import FindService from "../../components/Service/FindService";

export default function Explore() {
    return (
        <div>
            <HeaderBar activeBackground={true}/>
            <h1>Explore</h1>
            <FindService dark={false} />
        </div>
    );
}
