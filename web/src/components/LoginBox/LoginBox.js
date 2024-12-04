/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** LoginBox
*/

function LoginTexts() {
    return (
        <div className="text-center">
            <p className="text-white text-6xl font-spartan font-bold">
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
            <input type={id} id={id} class="bg-gray-500 border border-gray-700 text-white text-xl w-4/5 rounded-full focus:ring-blue-500 focus:border-blue-500 block p-4" placeholder={text} required />
        </div>
    )
}

function LoginTextFieldsBox( {text1, text2} ) {
    return (
        <div className="pt-10">
            <LoginTextField text={text1} id="email" />
            <LoginTextField text={text2} id="password"/>
        </div>
    )
}

function Button( {text} ) {
    return(
        <div className="flex justify-center pt-10">
            <button class="bg-white hover:bg-gray-300 text-black text-lg md:text-xl lg:text-2xl font-bold py-3 px-10 rounded-full">
                {text}
            </button>
        </div>
    )
}

export default function LoginBox () {
    return (
        <div className="bg-gradient-to-b from-zinc-700 to-gray-800 flex flex-col justify-center w-1/2 h-1/2 rounded-md">
            <LoginTexts />
            <LoginTextFieldsBox text1="Email" text2="Password" />
            <div className="text-center pt-10 text-white">
                You already have an account ? <a href="" className="font-bold text-white dark:text-white hover:underline">Register here !</a>
            </div>
           <Button text="Login" />    
        </div>
    )
}