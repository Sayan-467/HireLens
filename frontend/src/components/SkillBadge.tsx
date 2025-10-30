'use client';

interface SkillBadgeProps {
  skill: string;
  type: 'matching' | 'missing';
}

export default function SkillBadge({ skill, type }: SkillBadgeProps) {
  const isMatching = type === 'matching';

  return (
    <span
      className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium transition-all duration-200 ${
        isMatching
          ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 border border-green-300 dark:border-green-700'
          : 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300 border border-red-300 dark:border-red-700'
      }`}
    >
      <span
        className={`w-2 h-2 rounded-full mr-2 ${
          isMatching ? 'bg-green-500' : 'bg-red-500'
        }`}
      />
      {skill}
    </span>
  );
}
