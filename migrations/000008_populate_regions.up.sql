INSERT INTO regions (name)
VALUES ('Casablanca-Settat'), ('Souss-Massa'), ('Fès-Meknès'), ('Tanger-Tétouan-Al Hoceïma'), ('L''Oriental'), ('Rabat-Salé-Kénitra'), ('Béni Mellal-Khénifra'), ('Marrakech-Safi'), ('Drâa-Tafilalet'), ('Guelmim-Oued Noun'), ('Laâyoune-Sakia El Hamra'), ('Dakhla-Oued Ed-Dahab')
ON CONFLICT (name) DO NOTHING
