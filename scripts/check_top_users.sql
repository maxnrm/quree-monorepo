SELECT 
    users.chat_id, 
    COUNT(user_event_visits.id) AS VisitsCount
FROM 
    users 
LEFT JOIN 
    User_Event_Visits ON users.chat_id = user_event_visits.user_chat_id
GROUP BY 
    users.chat_id
ORDER BY 
    VisitsCount DESC;