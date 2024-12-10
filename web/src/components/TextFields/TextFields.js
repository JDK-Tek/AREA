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
            <input type={id} id={id} className="bg-gray-500 border border-gray-700 text-white 
                                          text-lg sm:text-xl md:text-2xl 
                                          w-11/12 sm:w-4/5 md:w-3/4 lg:w-2/3 
                                          rounded-full focus:ring-blue-500 focus:border-blue-500 block p-3 sm:p-4" 
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