import React from "react";

// const selectOptions = [
//     'All',
//     'By Votes',
//     'By Posted At',
//     'Users',
// ];

export default function FilterSearch({setNewQuery}) {
    if (!setNewQuery) {
        setNewQuery = () => {};
    }
    const [text, setText] = React.useState('');
    const handleClick = () => {
        setNewQuery(text);
    };
    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            handleClick();
        }
    };

    return (
        <div className='input-field flex justify-center my-[30px] gap-2'>
            <input type='text' placeholder='Search' className='p-2 rounded-md border-2 border-gray-300 w-1/2' onKeyDown={handleKeyDown} onInput={(e) => setText(e.target.value)} />
            {/* <select type='select' className='p-2 #3b82f6 hover:bg-[#E0E0E0]'>
                {
                    selectOptions.map((option, index) =>
                        <option key={index}>{option}</option>
                    )
                }
            </select> */}
            <button className='bg-[#860038] text-white hover:bg-[#680018] p-2 transition-colors duration-300'
                onClick={handleClick}
            >Search</button>
        </div>
    );
}