import React from "react";

export default function Topics() {
    return (
        <div className='topics-page w-full h-auto'>
            <div className='justify-center m-auto flex flex-col gap-3 mt-3'>
                <p className='text-5xl text-center'>Topics</p>
                <input
                    type='text'
                    className='w-1/2 m-auto p-2 rounded-md border-2 border-gray-300'
                    placeholder='Search for topics...'
                    onChange={(e) => console.log(e.target.value)} // Add your desired event handler here
                />
                <ul className='list-disc list-inside text-center'>
                    <li><a className='underline' href='#'>IEEE_CC</a></li>
                    <li><a className='underline' href='#'>Google_Toros</a></li>
                    <li><a className='underline' href='#'>Dr_Izaddoost_Club</a></li>
                </ul>
            </div>
        </div>
    );
}
