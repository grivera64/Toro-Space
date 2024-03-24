import React from 'react';
import { useNavigate } from 'react-router-dom';

import Discussions, {Discussion} from '../components/discussion/Discussions';
import FilterSearch from '../components/common/FilterSearch';
import Tabs from '../components/common/Tabs';

const tabs = [
    { label: 'Posts', content: 'Content 1' },
    { label: 'Discussions', content: 'Content 2' },
];

export default function Home() {
    const [activeTabIndex, setActiveTabIndex] = React.useState(0);
    const navigate = useNavigate();

    const handleSelect = (index) => {
        setActiveTabIndex(index);
        navigate(index === 0 ? '/posts' : '/discussions');
    };

    return (
        <div className='home-page w-full h-auto'>
            {/* <h1 className='text-4xl font-bold text-center'>Welcome to Toro Space</h1>
            <p className='text-lg text-center'>This is a simple web app that I built to learn about React and Tailwind CSS.</p>
            <br /> */}
            <br />
            <br />
            <Tabs tabs={tabs} selectedIndex={1} onSelect={handleSelect} />
            <FilterSearch />
            <Discussions>
                <Discussion />
                <Discussion />
                <Discussion />
            </Discussions>
        </div>
    )
}
