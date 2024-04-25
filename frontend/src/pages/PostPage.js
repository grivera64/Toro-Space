import React from 'react';
import {useParams} from 'react-router-dom';
import Posts, { Post } from '../components/post/Posts';

export default function PostPage() {
    const {postId} = useParams();
    const [post, setPost] = React.useState(null);

    React.useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch(`http://localhost:3030/posts/${postId}`, {
                    credentials: 'include'
                });
                const data = await response.json();
                setPost(data);
            } catch (error) {
                console.error('Failed to fetch post:', error);
                setPost({"error": error.message})
            }
        };
        fetchData();
    }, [postId]);

    if (post == null) {
        return (
            <div>
                <p>Loading...</p>
            </div>
        )
    }

    if (post.error) {
        return (
            <div>
                <p>Error: {post.error}</p>
            </div>
        )
    }

    return (
        <div className='p-10'>
            <button className='' onClick={() => {
                if (window.history.length > 2) {
                    window.location.href = document.referrer;
                } else {
                    window.location.href = '/home';
                }
            }}>Back</button>
            <br />
            <Posts>
                <Post postData={post} showLink={false} />
            </Posts>
        </div>
    );
}
