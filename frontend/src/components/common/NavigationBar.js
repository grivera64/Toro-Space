import React from 'react';
import { Link } from 'react-router-dom';

const navItems = [
    "Home",
    "Topics",
    "Organizations",
];

export default function NavigationBar() {
    return (
        <nav className='navigation-bar bg-[#860038] flex justify-between pr-20 pl-20 py-2'>
            <h1 className='text-3xl font-bold text-white'>Toro Space</h1>
            <ul className='flex items-center gap-[4vw]'>
                {
                    navItems.map((item, index) =>
                        <li
                            key={index}
                            className='hover:underline hover:cursor-pointer text-white text-lg font-bold transition duration-300 ease-in-out'
                        >
                            <Link to='/'>{item}</Link>
                        </li>
                    )
                }
                <button
                    className='bg-[#E6BC46] text-white text-lg font-bold py-2 px-4 rounded-full hover:bg-[#C69C26] hover:text-gray transition duration-300 ease-in-out'
                >
                    Sign In
                </button>
            </ul>
        </nav>
    );
}
