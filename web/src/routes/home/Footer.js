/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Footer
*/

import LinkListKit from '../../components/Link/LinkListKit';
import Download from '../../components/Download';

import Logo from './../../assets/fullLogo.png';

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
            <div className="mt-5 flex flex-wrap gap-7 p-5">
                <LinkListKit listLinks={dataLinkLists}/>
                <Download />
            </div>
        </div>
    );
}
