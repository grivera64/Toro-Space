import React, {useContext} from 'react';
import { Link } from 'react-router-dom';

// import Tabs from '../components/common/Tabs';
import { UserContext } from '../contexts/userContext';

import PostsView from '../views/PostsView';
// import DiscussionsView from '../views/DiscussionsView';

const tabs = [
    { label: 'Posts', content: 'Content 1' },
    { label: 'Discussions', content: 'Content 2' },
];

export default function Home() {
    const {user} = useContext(UserContext);
    const [tabIndex, setTabIndex] = React.useState(0);

    const handleSelect = (index) => {
        setTabIndex(index)
    };

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
            {/* <Tabs tabs={tabs} onSelect={handleSelect} /> */}
            {/*
                tabIndex === 0 &&*/
                <PostsView />
            }
            {/* {
                tabIndex === 1 &&
                <DiscussionsView />
            } */}
        </div>
    );
}
