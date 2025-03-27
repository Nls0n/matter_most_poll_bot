if not box.space.polls then
    box.schema.create_space('polls', {
        format = {
            {name = 'id', type = 'string'},        -- ID голосования
            {name = 'creator_id', type = 'string'}, -- ID создателя
            {name = 'question', type = 'string'},   -- Вопрос голосования
            {name = 'options', type = 'array'},     -- Варианты ответов
            {name = 'votes', type = 'map'}          -- Результаты (вариант -> количество)
        }
    })

    box.space.polls:create_index('primary', {
        type = 'hash',
        parts = {'id'}
    })

    box.space.polls:create_index('creator_idx', {
        type = 'tree',
        parts = {'creator_id'},
        if_not_exists = true
    })

    box.space.polls:create_index('question_idx', {
        type = 'tree',
        parts = {'question'},
        if_not_exists = true
    })
end

if not box.schema.func.exists('update_vote') then
    box.schema.function.create('update_vote', {
        body = [[
            function(poll_id, option)
                local polls = box.space.polls
                local poll = polls:get(poll_id)
                if not poll then return nil end

                local votes = poll[5] or {}
                votes[option] = (votes[option] or 0) + 1

                polls:update(poll_id, {{'=', 5, votes}})
                return votes
            end
        ]],
        if_not_exists = true
    })
end

return true