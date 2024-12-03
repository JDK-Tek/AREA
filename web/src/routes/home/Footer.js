/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Footer
*/

import Logo from './../../assets/fullLogo.png';
import LinkListKit from '../../components/Link/LinkListKit';

import exploreLinkData from "./../../data/ExploreLinkData";
import topIntegrationsLinkData from "./../../data/TopIntegrationLinkData";
import latestStoriesData from '../../data/LatestStoriesLinkData';

const dataLinkLists = [
    exploreLinkData,
    topIntegrationsLinkData,
    latestStoriesData
];

export default function Footer() {

    return (
        <div className="bg-chartgray-200 p-14">
            <img className="w-[150px]" src={Logo} alt="Logo"/>
            <LinkListKit listLinks={dataLinkLists}/>
        </div>
    );
}
