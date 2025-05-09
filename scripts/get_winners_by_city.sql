SELECT 
    users.chat_id, 
    users.quiz_city_name,
    users.date_quiz_finished,
    COUNT(user_event_visits.id) AS VisitsCount
FROM 
    users 
LEFT JOIN 
    User_Event_Visits ON users.chat_id = user_event_visits.user_chat_id
WHERE
    users.quiz_city_name is not null AND
    users.date_quiz_finished >= '2024-03-15'
GROUP BY 
    users.chat_id, users.quiz_city_name, users.date_quiz_finished
HAVING 
    COUNT(user_event_visits.id) >= 4
ORDER BY 
    users.quiz_city_name DESC, users.date_quiz_finished ASC;