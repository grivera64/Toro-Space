import React from 'react';
import { useNavigate, Link } from 'react-router-dom';

import Posts, {Post} from '../components/post/Posts';
import FilterSearch from '../components/common/FilterSearch';
import Tabs from '../components/common/Tabs';

const tabs = [
    { label: 'Posts', content: 'Content 1' },
    { label: 'Discussions', content: 'Content 2' },
];

export default function Home() {
    const navigate = useNavigate();

    const handleSelect = (index) => {
        navigate(index === 0 ? '/posts' : '/discussions');
    };

    return (
        <div className='home-page w-full h-auto'>
            {/* <h1 className='text-4xl font-bold text-center'>Welcome to Toro Space</h1>
            <p className='text-lg text-center'>This is a simple web app that I built to learn about React and Tailwind CSS.</p>
            <br /> */}
            <br />
            <br />
            <Tabs tabs={tabs} onSelect={handleSelect} />
            <FilterSearch />
            <Posts>
                <Post username={'IEEE_CC'} topics={["Computer_Science"]} content={<p>Come join our club: <a className='underline' href='https://torolink.csudh.edu/organization/ieee'>https://torolink.csudh.edu/organization/ieee</a></p>} />
                <Post username={'Google_Toros'} topics={["Computer_Science", "Careers"]} content={<p>Wish we were a club? Make it a reality and become President of our club!</p>} />
                <Post username={'Dr_Izaddoost_Club'} topics={["Careers"]} content={<p>Good luck on your presentations!</p>}/>
            </Posts>
        </div>
    )
}
