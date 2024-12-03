/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** LinkListKit
*/

import LinkList from "./LinkList";

export default function LinkListKit({ listLinks }) {
    return (
        <div className="p-5">
            <div className="mt-5 flex flex-wrap gap-7 p-5">
                {listLinks.map((listLink, index) => (
                    <LinkList
                        key={index}
                        title={listLink.title}
                        links={listLink.links}
                    />
                ))}
            </div>
        </div>
    );
}
