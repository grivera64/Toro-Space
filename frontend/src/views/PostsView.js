import React from 'react';

import Posts, {Post} from '../components/post/Posts';
import FilterSearch from '../components/common/FilterSearch';

import { UserContext } from '../contexts/userContext';

export default function PostsView() {
    const [posts, setPosts] = React.useState(null);
    const [hasNextPage, setHasNextPage] = React.useState(false);
    const [hasPrevPage, setHasPrevPage] = React.useState(false);

    const {user} = React.useContext(UserContext);
    const [latestPost, setLatestPost] = React.useState(null);
    const [searchQuery, setSearchQuery] = React.useState('');

    const [refreshNeeded, startRefresh] = React.useState(null);

    const [endpoint, setEndpoint] = React.useState('/posts?pageSize=10')
    React.useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch(`http://localhost:3030${endpoint}&search_query=${searchQuery}`, {
                    credentials: 'include'
                });
                const data = await response.json();
                console.log(data);
                setPosts(data['posts']);
                setHasNextPage(data['has_before']);
                setHasPrevPage(data['has_after']);
            } catch (error) {
                setPosts([]);
                setHasNextPage(false);
                setHasPrevPage(false);
                console.error('Failed to fetch posts:', error);
            }
        }
        fetchData();
    }, [latestPost, searchQuery, endpoint, refreshNeeded]);
    const [newPostContent, setNewPostContent] = React.useState('');

    const handlePostClick = async () => {
        if (newPostContent == null || newPostContent.length === 0) {
            return;
        }

        const nonHashtagMessage = newPostContent.replace(/#[a-zA-Z0-9]+/g, '');
        const hashtags = [...new Set(newPostContent.match(/#[a-zA-Z0-9]+/g) || [])];
        const response = await fetch(`http://localhost:3030/account/self/user/${user.id}/post`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            // Remove the hashtags from the content and send them as topics
            body: JSON.stringify({ content: nonHashtagMessage, topics: hashtags.map(tag => tag.substring(1))})
        });

        if (response.status !== 200) {
            console.error('Failed to create post');
            return;
        }
        const data = await response.json();
        setLatestPost(data);
        setNewPostContent('');
    };

    if (posts == null || user == null) {
        return (
            <div className='p-10'>
                <FilterSearch setSearchQuery={setSearchQuery} />
                <p>Loading...</p>
            </div>
        )
    }

    return (
        <div className='p-10'>
            <FilterSearch setSearchQuery={setSearchQuery} />
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
                        className="ml-2 px-4 py-2 bg-blue-500
                            hover:bg-[#1b62d6] text-white rounded-md transition-colors duration-300"
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
                {posts && posts.map(post => (
                    <Post
                        key={post.id}
                        // postID={post.id}
                        // displayName={post.author['display_name']}
                        // avatar={post.author['avatar_url']}
                        // topics={post.topics}
                        // date={post['created_at']}
                        // likes={post.likes}
                        // content={post.content}
                        // isLiked={post.liked_by.map(obj => obj.id).some(id => id === user.id)} // Check if current user ID is in the likedBy list
                        postData={post}
                        showLink={true}
                        startRefresh={startRefresh}
                    />
                ))}
            </Posts>
            <div className="flex justify-center mt-4 gap-2">
                <button
                    className="px-4 py-2 bg-[#860038] 
                        hover:bg-[#680018]
                        disabled:bg-gray-500 disabled:hover:cursor-not-allowed
                        text-white rounded-md transition-colors duration-300"
                    onClick={() => setEndpoint(`/posts?pageSize=10&after=${posts[0].id}`)}
                    disabled={!hasPrevPage}
                >
                    Previous Page
                </button>
                <button
                    className="px-4 py-2 bg-[#860038]
                        hover:bg-[#680018]
                        disabled:bg-gray-500 disabled:hover:cursor-not-allowed
                        text-white rounded-md transition-colors duration-300"
                    onClick={() => setEndpoint(`/posts?pageSize=10&before=${posts[posts.length - 1].id}`)}
                    disabled={!hasNextPage}
                >
                    Next Page
                </button>
            </div>
        </div>
    );
}