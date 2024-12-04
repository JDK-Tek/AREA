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

const dataFooterLinks = [
    {
        title: "Developers",
        url: "/developers"
    },
    {
        title: "About Us",
        url: "/about"
    },
    {
        title: "Contact Us",
        url: "/contact"
    },
    {
        title: "Privacy Policy",
        url: "/privacy"
    },
    {
        title: "Terms of Service",
        url: "/terms"
    }
]

function FooterLink()
{

    return (
        <div className="flex flex-wrap gap-7 pl-[120px] pt-10">
            {dataFooterLinks.map((link, index) => (
                <label
                    key={index}
                    className="cursor-pointer block mb-4 hover:text-chartgray-100 text-white font-bold"
                    onClick={() => { window.location.href = link.url }}
                >
                    {link.title}
                </label>
            ))}
        </div>
    );
}

export default function Footer() {

    return (
        <div className="bg-chartgray-200 p-14">
            <img className="w-[150px]" src={Logo} alt="Logo"/>
            <div className="mt-5 p-5">
                <div className="flex flex-wrap gap-7">
                    <LinkListKit listLinks={dataLinkLists}/>
                    <Download />
                </div>
                <FooterLink />
            </div>
        </div>
    );
}
