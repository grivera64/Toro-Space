import React from "react";

export default function Posts(props) {
    return (
        <div className='posts container space-y-2 flex flex-col justify-center items-center mx-auto w-auto'>    
            {props.children}
        </div>
    );
}

export function Post({displayName, avatar, content, date, topics}) {
    const [upvoteSelected, setUpvoteSelected] = React.useState(false);
    const [downvoteSelected, setDownvoteSelected] = React.useState(false);
    const [likes, setLikes] = React.useState(0);

    if (date !== undefined && date !== null && date !== '') {
        date = new Date(date).toLocaleString();
    } else {
        date = 'Unknown Date';
    }

    return (
        <div className='post container rounded-md border-2 border-gray-300 w-1/2'>
            <div className='post-header text-xs flex flex-row text-wrap gap-3 p-2 bg-[#DDDDDD] w-full rounded-t'>
                <img src={avatar} alt='avatar' />
                <p>@{displayName ?? 'placeholder'}</p>
                <ul className='flex flex-row gap-2 text-blue-500'>
                    {
                        topics?.map((topic, index) => (
                            <li key={index}>topic/{topic}</li>
                        ))
                    }
                </ul>
                <p>{date}</p>
            </div>
            {/* <hr /> */}
            <div className='post-content p-3 text-lg'>
                <p>{content ?? 'This is a placeholder post.'}</p>
            </div>
            <div className='vote-footer p-3'>
                <div className='vote-footer-left flex flex-row gap-3'>
                    <button className={`vote-footer-left-like ${upvoteSelected ? 'bg-blue-200' : ''}`}
                        onClick={() => {
                            setLikes((likes) => {
                                if (downvoteSelected) {
                                    return likes + 2;
                                } else if (!upvoteSelected) {
                                    return likes + 1;
                                } else {
                                    return likes - 1;
                                }
                            });
                            setUpvoteSelected((upvoteSelected) => !upvoteSelected);
                            setDownvoteSelected(false);
                        }}
                    >
                        {/* <img width={30} height={30} src='/images/upvote.svg'/> */}
                        <p>👍</p>
                    </button>
                    <p>{likes ?? 0}</p>
                    <button className={`vote-footer-left-dislike ${downvoteSelected ? 'bg-red-200' : ''}`}
                        onClick={() => {
                            setLikes((likes) => {
                                if (upvoteSelected) {
                                    return likes - 2;
                                } else if (!downvoteSelected) {
                                    return likes - 1;
                                } else {
                                    return likes + 1;
                                }
                            });
                            setDownvoteSelected((downvoteSelected) => !downvoteSelected);
                            setUpvoteSelected(false);
                        }}
                    >
                    {/* <img width={30} height={30} src='/images/downvote.svg'/> */}
                        <p>👎</p>
                    </button>
                    <p>💬</p>
                </div>
            </div>
        </div>
    );
}
