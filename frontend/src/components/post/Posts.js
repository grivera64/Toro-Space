import React from "react";
import { UserContext } from "../../contexts/userContext";

export default function Posts(props) {
    return (
        <div className='posts container space-y-2 flex flex-col justify-center items-center mx-auto w-auto'>    
            {props.children}
        </div>
    );
}

// export function Post({postID, displayName, avatar, content, date, topics, likes, isLiked}) {
export function Post({postData: {id, author: {display_name, avatar_url, role}, content, created_at, topics, likes, liked_by}}) {
    const {user} = React.useContext(UserContext);
    const [upvoteSelected, setUpvoteSelected] = React.useState(liked_by.map(u => u.id).some(id => id === user.id));
    const [likesCount, setLikesCount] = React.useState(likes);

    if (created_at !== undefined && created_at !== null && created_at !== '') {
        created_at = new Date(created_at).toLocaleString();
    } else {
        created_at = 'Unknown Date';
    }

    React.useEffect(() => {
        const isLikedByUser = liked_by.some(u => u.id === user.id);
        setUpvoteSelected(isLikedByUser);
    }, [liked_by, user.id]);

    return (
        <div className='post container rounded-md border-2 border-gray-300 w-1/2'>
            <div className='post-header text-xs flex flex-row text-wrap gap-3 p-2 bg-[#DDDDDD] w-full rounded-t'>
                <img src={avatar_url} alt='avatar' />
                <p>@{display_name ?? 'placeholder'}</p>
                <ul className='flex flex-row gap-2 text-blue-500'>
                    {
                        topics?.map((topic, index) => (
                            <li key={index}>topic/{topic}</li>
                        ))
                    }
                </ul>
                <p>{created_at}</p>
            </div>
            {/* <hr /> */}
            <div className='post-content p-3 text-lg'>
                <p>{content ?? 'This is a placeholder post.'}</p>
            </div>
            <div className='vote-footer p-3'>
                <div className='vote-footer-left flex flex-row gap-3'>
                    <button className={`vote-footer-left-like ${upvoteSelected ? 'bg-blue-200' : ''}`}
                        onClick={() => {
                            if (upvoteSelected) {
                                fetch(`http://localhost:3030/posts/${id}/like/?type=unlike`, {
                                    method: 'POST', credentials: 'include'
                                }).then(response => response.json())
                                    .then(data => setLikesCount(data.likes))
                                    .catch(error => console.log(error));
                                setUpvoteSelected(false);
                            } else {
                                fetch(`http://localhost:3030/posts/${id}/like/?type=like`, {
                                    method: 'POST', credentials: 'include'
                                }).then(response => response.json())
                                    .then(data => setLikesCount(data.likes))
                                    .catch(error => console.log(error));
                                setUpvoteSelected(true);
                            }
                        }}
                    >
                        {/* <img width={30} height={30} src='/images/upvote.svg'/> */}
                        <p>👍</p>
                    </button>
                    <p>{likesCount ?? 0}</p>
                    {/* <p>💬</p> */}
                </div>
            </div>
        </div>
    );
}
