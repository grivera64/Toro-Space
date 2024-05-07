import React, {useContext} from 'react';
import { Link } from 'react-router-dom';

import { UserContext } from '../contexts/userContext';

import PostsView from '../views/PostsView';

export default function Home() {
    const {user} = useContext(UserContext);

    return (
        <div className='home-page w-full h-auto'>
            <div className='flex justify-center my-4'>
                {
                    user?.error ||
                    <div className='flex flex-col gap-5'>
                        <p className='text-center'>Logged in as {user['display_name']}</p>
                        <p className='text-center'><Link to='/select' className='cursor-pointer hover:underline'>Switch Accounts</Link></p>
                    </div>
                    
                }
            </div>
            <PostsView />
        </div>
    );
}
