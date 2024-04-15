import React from 'react';

import Posts, {Post} from '../components/post/Posts';
import FilterSearch from '../components/common/FilterSearch';

import { UserContext } from '../contexts/userContext';

export default function PostsView() {
    const [posts, setPosts] = React.useState([]);
    const {user} = React.useContext(UserContext);

    React.useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch('http://localhost:3030/posts/?pageSize=10');
                const data = await response.json();
                console.log(data);
                setPosts(data);
            } catch (error) {
                console.error('Failed to fetch posts:', error);
            }
        }
        
        fetchData();
    }, []);
    const [newPostContent, setNewPostContent] = React.useState('');

    const handlePostClick = async () => {
        const response = await fetch(`http://localhost:3030/account/self/user/${user.id}/post`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ content: newPostContent })
        });

        if (response.status !== 200) {
            console.error('Failed to create post');
            return;
        }
        const data = await response.json();
        console.log(data);
        // Update the posts state with the new post
        setPosts([data, ...posts]);
        // Clear the new post content
        setNewPostContent('');
    };

    return (
        <div className='p-10'>
            <FilterSearch />
            {
                user.role === 'organization' &&
                <div className="flex justify-center items-center m-20">
                    <textarea
                        className="w-full h-24 p-2 border border-gray-300 rounded-md resize-none"
                        placeholder="What's your next greatest event?"
                        value={newPostContent}
                        onChange={(e) => setNewPostContent(e.target.value)}
                    ></textarea>
                    <button
                        className="ml-2 px-4 py-2 bg-blue-500 text-white rounded-md"
                        onClick={handlePostClick}
                    >
                        Post
                    </button>
                </div>
            }
            <Posts>
                {/* <Post username={'IEEE_CC'} topics={["Computer_Science"]} date={Date.now()} content={<p>Come join our club: <a className='underline' href='https://torolink.csudh.edu/organization/ieee'>https://torolink.csudh.edu/organization/ieee</a></p>} />
                <Post username={'Google_Toros'} topics={["Computer_Science", "Careers"]} content={<p>Wish we were a club? Make it a reality and become President of our club!</p>} />
                <Post username={'Dr_Izaddoost_Club'} topics={["Careers"]} content={<p>Good luck on your presentations!</p>} /> */}
                {posts && posts.map(post => <Post key={post.id} displayName={post.author['display_name']} avatar={post.author['avatar_url']} topics={post.topics} date={post['created_at']} content={post.content} />)}
            </Posts>
        </div>
    );
}