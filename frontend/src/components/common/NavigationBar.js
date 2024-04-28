import React, {useContext} from 'react';
import { Link } from 'react-router-dom';

import { UserContext } from '../../contexts/userContext';

const navItems = [
    {
        name: "Home",
        path: "/",
    },
    {
        name: "Topics",
        path: "/topics",
    },
    {
        name: "Organizations",
        path: "/organizations"
    },
];

export default function NavigationBar() {
    const {user, loggedIn} = useContext(UserContext);
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
                            <Link to={item.path}>{item.name}</Link>
                        </li>
                    )
                }
                {
                    user?.role === 'admin' && <li className='hover:underline hover:cursor-pointer text-white text-lg font-bold transition duration-300 ease-in-out'><Link to='/admin'>Admin</Link></li>
                }
                {
                    (!loggedIn) &&
                    <button
                    className='bg-[#E6BC46] text-white text-lg font-bold py-2 px-4 rounded-full hover:bg-[#C69C26] hover:text-gray transition duration-300 ease-in-out'
                    onClick={() => {
                        window.location.href = 'http://localhost:3030/auth/google';
                    }}
                    >
                        Sign In
                    </button>
                }
                {
                    (loggedIn) &&
                    <button
                    className='bg-[#E6BC46] text-white text-lg font-bold py-2 px-4 rounded-full hover:bg-[#C69C26] hover:text-gray transition duration-300 ease-in-out'
                    onClick={() => {
                        window.location.href = 'http://localhost:3030/logout';
                    }}
                    >
                        Sign Out
                    </button>
                }
            </ul>
        </nav>
    );
}
