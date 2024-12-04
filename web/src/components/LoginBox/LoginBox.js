/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** LoginBox
*/

function LoginTexts() {
    return (
        <div className="text-center">
            <p className="text-white text-6xl pt-10 font-spartan font-bold">
                LOGIN
            </p>
            <p className="text-violet-600 font-bold text-2xl">
                Nice to see you again
            </p>
        </div>
    )
}

function LoginTextField ( {text, id} ) {
    return (
        <div className="pt-5 justify-center flex">
            <input type={id} id={id} class="bg-gray-50 border border-gray-700 text-gray-900 text-lg w-4/5 rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5" placeholder={text} required />
        </div>
    )
}

function LoginTextFieldsBox( {text1, text2}) {
    return (
        <div className="pt-10">
            <LoginTextField text={text1} id="email" />
            <LoginTextField text={text2} id="password"/>
        </div>
    )
}

export default function LoginBox () {
    return (
        <div className="bg-gradient-to-b from-zinc-700 to-gray-800 flex flex-col w-1/2 h-3/4 rounded-md">
            <LoginTexts />
            <LoginTextFieldsBox text1="Email" text2="Password" />
        </div>
    )
}