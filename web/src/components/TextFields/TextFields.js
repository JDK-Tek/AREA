/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** TextFields
*/

// Login / Register TextFields

export function LRTextField ( {text, id, handleChangeField} ) {
    return (
        <div className="pt-5 justify-center flex">
            <input type={id} id={id} className="bg-chartgray-300 border border-chartpurple-200 text-white 
                                          lg:text-lg md:text-md text-sm
                                          w-11/12 sm:w-4/5 md:w-3/4 lg:w-2/3 
                                          rounded-full focus:outline-none focus:ring-2 focus:ring-chartpurple-100 block p-3 sm:p-4" 
                   placeholder={text} required onChange={handleChangeField}/>
        </div>
    )
}

export function LRTextFieldsBox( {text1, text2, handleChangeField} ) {
    return (
        <div className="pt-10">
            <LRTextField text={text1} id="email" handleChangeField={handleChangeField}/>
            <LRTextField text={text2} id="password" handleChangeField={handleChangeField}/>
        </div>
    )
}