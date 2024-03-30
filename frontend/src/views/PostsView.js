import React from 'react';

import Posts, {Post} from '../components/post/Posts';
import FilterSearch from '../components/common/FilterSearch';

export default function PostsView() {
    return (
        <div>
            <FilterSearch />
            <Posts>
                <Post username={'IEEE_CC'} topics={["Computer_Science"]} date={Date.now()} content={<p>Come join our club: <a className='underline' href='https://torolink.csudh.edu/organization/ieee'>https://torolink.csudh.edu/organization/ieee</a></p>} />
                <Post username={'Google_Toros'} topics={["Computer_Science", "Careers"]} content={<p>Wish we were a club? Make it a reality and become President of our club!</p>} />
                <Post username={'Dr_Izaddoost_Club'} topics={["Careers"]} content={<p>Good luck on your presentations!</p>} />
            </Posts>
        </div>
    );
}