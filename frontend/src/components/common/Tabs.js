import React, { useState } from 'react';

function Tab({ label, isActive, onSelect }) {
  const activeClassName = isActive ? 'bg-[#F1F1F1]' : '';

  return (
    <li className={`text-center py-2 px-4 rounded-md cursor-pointer ${activeClassName}`} onClick={onSelect}>
      {label}
    </li>
  );
}

function TabContent({ children, isActive }) {
  return isActive ? <div className="mt-4">{children}</div> : null;
}

function Tabs({ tabs, selectedIndex = 0, onSelect }) {
  const [activeIndex, setActiveIndex] = useState(selectedIndex);

  const handleClick = (index) => {
    setActiveIndex(index);
    onSelect && onSelect(index);
  };

  return (
    <div className="text-center">
      <ul className='flex flex-row gap-2 justify-center'>
        {tabs.map((tab, index) => (
          <Tab
            key={index}
            label={tab.label}
            isActive={index === activeIndex}
            onSelect={() => handleClick(index)}
          />
        ))}
      </ul>
    </div>
  );
}

export default Tabs;
